package app

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	db "github.com/cometbft/cometbft-db"
)

var (
	stateKey = []byte("stateKey")
)

type State struct {
	db *db.GoLevelDB

	// Size is essentially the amount of transactions that have been processed.
	// This is used for the appHash
	Size   int64 `json:"size"`
	Height int64 `json:"height"`
}

// NewState creates a new State instance with the given path.
// Returns an error if the state cannot be created or loaded.
func NewState(path string) (*State, error) {
	// Ensure the directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory at %s: %w", path, err)
	}

	name := "state"
	db, err := db.NewGoLevelDB(name, path)
	if err != nil {
		return nil, fmt.Errorf("failed to create persistent state at %s: %w", path, err)
	}

	state := &State{db: db}
	stateBytes, err := db.Get(stateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	if len(stateBytes) > 0 {
		if err := json.Unmarshal(stateBytes, state); err != nil {
			return nil, fmt.Errorf("failed to read current state: %w", err)
		}
	}

	return state, nil
}

// Save persists the current state to the database.
func (s *State) Save() error {
	stateBytes, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := s.db.Set(stateKey, stateBytes); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// Hash returns a byte slice representing the state's hash.
func (s *State) Hash() []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(s.Size))
	return bytes
}

// Close closes the state's database connection.
func (s *State) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close state database: %w", err)
	}
	return nil
}
