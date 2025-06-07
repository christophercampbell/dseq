package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/christophercampbell/dseq/app"
	"github.com/christophercampbell/dseq/app/config"
	"github.com/cometbft/cometbft/libs/log"
)

// TestConfig creates a test configuration
func TestConfig(t *testing.T) *config.Config {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "dseq-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return config.NewConfig(tmpDir)
}

// TestState creates a test state
func TestState(t *testing.T) *app.State {
	t.Helper()

	cfg := TestConfig(t)
	state, err := app.NewState(cfg.HomeDir)
	if err != nil {
		t.Fatalf("failed to create test state: %v", err)
	}

	t.Cleanup(func() {
		if err := state.Close(); err != nil {
			t.Errorf("failed to close test state: %v", err)
		}
	})

	return state
}

// TestLogger creates a test logger
func TestLogger(t *testing.T) log.Logger {
	t.Helper()
	return log.NewNopLogger()
}

// TestSequencer creates a test sequencer
func TestSequencer(t *testing.T) *app.SequencerApplication {
	t.Helper()

	state := TestState(t)
	logger := TestLogger(t)

	sequencer, err := app.NewSequencer(
		logger,
		app.WithIdentity("test-sequencer"),
		app.WithState(state),
	)
	if err != nil {
		t.Fatalf("failed to create test sequencer: %v", err)
	}

	return sequencer
}

// CreateTestFiles creates test files in the given directory
func CreateTestFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()

	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("failed to create directory for %s: %v", name, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", name, err)
		}
	}
}
