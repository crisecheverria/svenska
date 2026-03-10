package updater

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/creativeprojects/go-selfupdate"
)

const repo = "crisecheverria/svenska"

type UpdateInfo struct {
	Available  bool
	NewVersion string
}

// CheckForUpdate checks GitHub releases for a newer version (non-blocking friendly).
func CheckForUpdate(currentVersion string) UpdateInfo {
	if currentVersion == "dev" || currentVersion == "" {
		return UpdateInfo{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return UpdateInfo{}
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{Source: source})
	if err != nil {
		return UpdateInfo{}
	}

	latest, found, err := updater.DetectLatest(ctx, selfupdate.ParseSlug(repo))
	if err != nil || !found {
		return UpdateInfo{}
	}

	if latest.GreaterThan(currentVersion) {
		return UpdateInfo{Available: true, NewVersion: latest.Version()}
	}
	return UpdateInfo{}
}

// DoUpdate downloads and applies the latest release.
func DoUpdate(currentVersion string) error {
	if currentVersion == "dev" || currentVersion == "" {
		return fmt.Errorf("cannot update a development build — install from a release")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return fmt.Errorf("failed to create update source: %w", err)
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{Source: source})
	if err != nil {
		return fmt.Errorf("failed to create updater: %w", err)
	}

	latest, found, err := updater.DetectLatest(ctx, selfupdate.ParseSlug(repo))
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	if !found {
		return fmt.Errorf("no release found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	if !latest.GreaterThan(currentVersion) {
		return fmt.Errorf("already up to date (v%s)", currentVersion)
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}

	if err := updater.UpdateTo(ctx, latest, exe); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("Updated to v%s\n", latest.Version())
	return nil
}
