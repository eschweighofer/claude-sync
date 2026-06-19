package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/eschweighofer/claude-sync/internal/config"
)

// PullFromLocalWSL copies VS Code extension data from a Windows-mounted source
// path into the local WSL VS Code user-data directory. This mode is local-only
// and does not use any external provider.
func PullFromLocalWSL(cfg *config.Config, dryRun bool) (*SyncResult, error) {
	if cfg.LocalWSLSource == "" {
		return nil, fmt.Errorf("local WSL source is not configured")
	}

	source := filepath.Clean(cfg.LocalWSLSource)
	target, err := localWSLTargetUserData(cfg.LocalWSLProfile)
	target = filepath.Clean(target)
	if target == "" {
		return nil, fmt.Errorf("could not resolve local WSL VS Code target directory")
	}

	if info, err := os.Stat(source); err != nil || !info.IsDir() {
		if err != nil {
			return nil, fmt.Errorf("local WSL source not accessible: %w", err)
		}
		return nil, fmt.Errorf("local WSL source is not a directory: %s", source)
	}

	files, err := GetLocalFiles(source, vscodeSyncPaths)
	if err != nil {
		return nil, err
	}

	result := &SyncResult{}
	keys := make([]string, 0, len(files))
	for rel := range files {
		keys = append(keys, rel)
	}
	sort.Strings(keys)

	for _, relPath := range keys {
		srcPath := filepath.Join(source, relPath)
		dstPath := filepath.Join(target, relPath)

		// Guard against traversal from crafted source trees
		cleanDst := filepath.Clean(dstPath)
		if !strings.HasPrefix(cleanDst, target+string(filepath.Separator)) && cleanDst != target {
			result.Errors = append(result.Errors, fmt.Errorf("refusing to write outside %s: %s", target, relPath))
			continue
		}

		if dryRun {
			result.Downloaded = append(result.Downloaded, relPath)
			continue
		}

		if err := copyIfChanged(srcPath, dstPath); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", relPath, err))
			continue
		}
		result.Downloaded = append(result.Downloaded, relPath)
	}

	return result, nil
}

func localWSLTargetUserData(profile string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.EqualFold(profile, "insiders") {
		return filepath.Join(home, ".vscode-server-insiders", "data"), nil
	}
	return filepath.Join(home, ".vscode-server", "data"), nil
}

func copyIfChanged(srcPath, dstPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if dstInfo, err := os.Stat(dstPath); err == nil {
		if dstInfo.Size() == srcInfo.Size() && sameSecond(dstInfo.ModTime(), srcInfo.ModTime()) {
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0700); err != nil {
		return err
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return err
	}
	if err := dst.Close(); err != nil {
		return err
	}

	return os.Chtimes(dstPath, srcInfo.ModTime(), srcInfo.ModTime())
}

func sameSecond(a, b time.Time) bool {
	return a.UTC().Truncate(time.Second).Equal(b.UTC().Truncate(time.Second))
}
