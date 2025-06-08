package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	// Create a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "state_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid path",
			path:    tmpDir,
			wantErr: false,
		},
		{
			name:    "creates directory if not exists",
			path:    filepath.Join(tmpDir, "nonexistent"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := NewState(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, state)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, state)
				assert.Equal(t, int64(0), state.Size)
				assert.Equal(t, int64(0), state.Height)

				// Verify directory was created
				if tt.name == "creates directory if not exists" {
					_, err := os.Stat(tt.path)
					assert.NoError(t, err, "directory should have been created")
				}
			}
		})
	}
}

func TestStateSaveAndLoad(t *testing.T) {
	// Create a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "state_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create initial state
	state, err := NewState(tmpDir)
	require.NoError(t, err)
	require.NotNil(t, state)

	// Modify state
	state.Size = 42
	state.Height = 100

	// Save state
	err = state.Save()
	require.NoError(t, err)

	// Close state
	err = state.Close()
	require.NoError(t, err)

	// Load state again
	newState, err := NewState(tmpDir)
	require.NoError(t, err)
	require.NotNil(t, newState)

	// Verify state was loaded correctly
	assert.Equal(t, int64(42), newState.Size)
	assert.Equal(t, int64(100), newState.Height)

	// Clean up
	err = newState.Close()
	require.NoError(t, err)
}

func TestStateHash(t *testing.T) {
	// Create a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "state_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create state
	state, err := NewState(tmpDir)
	require.NoError(t, err)
	require.NotNil(t, state)
	defer state.Close()

	// Test hash with different sizes
	tests := []struct {
		name     string
		size     int64
		expected []byte
	}{
		{
			name:     "zero size",
			size:     0,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "positive size",
			size:     42,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 42},
		},
		{
			name:     "large size",
			size:     1 << 32,
			expected: []byte{0, 0, 0, 1, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state.Size = tt.size
			hash := state.Hash()
			assert.Equal(t, tt.expected, hash)
		})
	}
}

func TestStateClose(t *testing.T) {
	// Create a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "state_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create state
	state, err := NewState(tmpDir)
	require.NoError(t, err)
	require.NotNil(t, state)

	// Close state
	err = state.Close()
	require.NoError(t, err)

	// Try to close again
	err = state.Close()
	require.Error(t, err)
}
