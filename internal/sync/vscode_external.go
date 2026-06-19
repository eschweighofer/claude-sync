package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const vscodeRemotePrefix = "_external/vscode-user-data/"

var vscodeSyncPaths = []string{
	"User/globalStorage",
	"User/workspaceStorage",
}

func (s *Syncer) PushVSCodeSource(ctx context.Context, sourceRoot string) (*SyncResult, error) {
	sourceRoot = filepath.Clean(sourceRoot)
	if sourceRoot == "" {
		return nil, fmt.Errorf("vscode sync source is not configured")
	}
	if fi, err := os.Stat(sourceRoot); err != nil || !fi.IsDir() {
		if err != nil {
			return nil, fmt.Errorf("vscode sync source not accessible: %w", err)
		}
		return nil, fmt.Errorf("vscode sync source is not a directory: %s", sourceRoot)
	}

	files, err := GetLocalFiles(sourceRoot, vscodeSyncPaths)
	if err != nil {
		return nil, err
	}

	result := &SyncResult{}
	keys := make([]string, 0, len(files))
	for rel := range files {
		keys = append(keys, rel)
	}
	sort.Strings(keys)

	for _, rel := range keys {
		full := filepath.Join(sourceRoot, rel)
		data, err := os.ReadFile(full)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}

		compressed, err := gzipCompress(data)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}
		encrypted, err := s.encryptor.Encrypt(compressed)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}

		remoteKey := vscodeRemotePrefix + filepath.ToSlash(rel) + ".age"
		if err := s.storage.Upload(ctx, remoteKey, encrypted); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}
		result.Uploaded = append(result.Uploaded, rel)
	}

	return result, nil
}

func (s *Syncer) PullVSCodeSource(ctx context.Context, targetRoot string, dryRun bool) (*SyncResult, error) {
	targetRoot = filepath.Clean(targetRoot)
	if targetRoot == "" {
		return nil, fmt.Errorf("vscode sync source is not configured")
	}

	remoteObjects, err := s.storage.List(ctx, vscodeRemotePrefix)
	if err != nil {
		return nil, err
	}

	result := &SyncResult{}
	for _, obj := range remoteObjects {
		if !strings.HasSuffix(obj.Key, ".age") {
			continue
		}
		rel := strings.TrimPrefix(obj.Key, vscodeRemotePrefix)
		rel = strings.TrimSuffix(rel, ".age")
		if rel == "" {
			continue
		}
		dstPath := filepath.Join(targetRoot, rel)
		cleanDst := filepath.Clean(dstPath)
		if !strings.HasPrefix(cleanDst, targetRoot+string(filepath.Separator)) && cleanDst != targetRoot {
			result.Errors = append(result.Errors, fmt.Errorf("refusing to write outside %s: %s", targetRoot, rel))
			continue
		}

		if dryRun {
			result.Downloaded = append(result.Downloaded, rel)
			continue
		}

		encrypted, err := s.storage.Download(ctx, obj.Key)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}
		data, err := s.encryptor.Decrypt(encrypted)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}
		if isGzipped(data) {
			data, err = gzipDecompress(data)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
				continue
			}
		}
		if err := os.MkdirAll(filepath.Dir(dstPath), 0700); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}
		if err := os.WriteFile(dstPath, data, 0600); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel, err))
			continue
		}
		result.Downloaded = append(result.Downloaded, rel)
	}

	return result, nil
}
