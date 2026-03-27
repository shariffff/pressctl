package updater

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pressctl/cli/internal/installer"
)

const (
	repo    = "shariffff/pressctl"
	apiBase = "https://api.github.com"
)

// Release holds relevant GitHub release info.
type Release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// LatestRelease fetches the latest release from GitHub.
func LatestRelease() (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", apiBase, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "pressctl-updater")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach GitHub: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}
	return &r, nil
}

// IsNewer reports whether remoteTag is strictly newer than localVersion.
// Both may carry an optional "v" prefix.
func IsNewer(localVersion, remoteTag string) bool {
	local := strings.TrimPrefix(localVersion, "v")
	remote := strings.TrimPrefix(remoteTag, "v")
	return semverGreater(remote, local)
}

// Install downloads remoteTag and replaces the running binary and the
// ansible directory at ~/.pressctl/ansible/. The user config is never touched.
func Install(version string) error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	versionNum := strings.TrimPrefix(version, "v")
	archiveName := fmt.Sprintf("pressctl_%s_%s_%s", versionNum, goos, goarch)

	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	filename := archiveName + ext
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, version, filename)

	// Place temp dir inside ~/.pressctl/ to guarantee same-filesystem rename.
	pressctlDir := installer.GetWordmonDir()
	if err := os.MkdirAll(pressctlDir, 0755); err != nil {
		return fmt.Errorf("failed to prepare install dir: %w", err)
	}
	tmpDir, err := os.MkdirTemp(pressctlDir, "update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, filename)
	if err := downloadFile(url, archivePath); err != nil {
		return err
	}

	if err := extractTarGz(archivePath, tmpDir); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	extractedDir := filepath.Join(tmpDir, archiveName)

	if err := replaceBinary(extractedDir, goos); err != nil {
		return err
	}

	if err := replaceAnsible(filepath.Join(extractedDir, "ansible")); err != nil {
		return err
	}

	return nil
}

// semverGreater returns true if a > b using numeric segment comparison.
func semverGreater(a, b string) bool {
	pa := semverParts(a)
	pb := semverParts(b)
	for i := 0; i < 3; i++ {
		if pa[i] != pb[i] {
			return pa[i] > pb[i]
		}
	}
	return false
}

func semverParts(v string) [3]int {
	parts := strings.SplitN(v, ".", 3)
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	var out [3]int
	for i := 0; i < 3; i++ {
		out[i], _ = strconv.Atoi(parts[i])
	}
	return out
}

func downloadFile(url, dest string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "pressctl-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create download file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("failed to write download: %w", err)
	}
	return nil
}

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	cleanDest := filepath.Clean(dest)
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Zip-slip protection: clean the path before joining.
		target := filepath.Join(cleanDest, filepath.Clean("/"+hdr.Name)[1:])
		if target != cleanDest && !strings.HasPrefix(target, cleanDest+string(os.PathSeparator)) {
			return fmt.Errorf("invalid path in archive: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(out, tr)
			out.Close()
			if copyErr != nil {
				return copyErr
			}
		}
	}
	return nil
}

func replaceBinary(extractedDir, goos string) error {
	binaryName := "press"
	if goos == "windows" {
		binaryName = "press.exe"
	}
	newBinary := filepath.Join(extractedDir, binaryName)

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine current binary path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("could not resolve binary path: %w", err)
	}

	// Copy new binary to a sibling tmp file, then atomic rename.
	tmpBin := exePath + ".new"
	src, err := os.Open(newBinary)
	if err != nil {
		return fmt.Errorf("binary not found in archive: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(tmpBin, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to stage new binary: %w", err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(tmpBin)
		return fmt.Errorf("failed to copy new binary: %w", err)
	}
	dst.Close()

	if err := os.Rename(tmpBin, exePath); err != nil {
		os.Remove(tmpBin)
		return fmt.Errorf("failed to replace binary: %w", err)
	}
	return nil
}

func replaceAnsible(newAnsible string) error {
	if _, err := os.Stat(newAnsible); err != nil {
		return nil // no ansible dir in archive — skip
	}

	ansibleDir := installer.GetAnsibleDir()
	backupDir := ansibleDir + ".bak"

	// Remove any stale backup from a previous failed update.
	os.RemoveAll(backupDir)

	// Back up current ansible dir.
	if _, err := os.Stat(ansibleDir); err == nil {
		if err := os.Rename(ansibleDir, backupDir); err != nil {
			return fmt.Errorf("failed to backup ansible dir: %w", err)
		}
	}

	// Move new ansible into place (same FS → atomic).
	if err := os.Rename(newAnsible, ansibleDir); err != nil {
		os.Rename(backupDir, ansibleDir) // restore on failure
		return fmt.Errorf("failed to install ansible dir: %w", err)
	}

	os.RemoveAll(backupDir)
	return nil
}
