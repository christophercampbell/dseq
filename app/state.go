package app

import (
	"encoding/binary"
	"encoding/json"
	"log"

	"github.com/cometbft/cometbft-db"
)

var (
	stateKey = []byte("stateKey")
)

type State struct {
	db *db.GoLevelDB

	// Size is essentially the amount of transactions that have been processes.
	// This is used for the appHash
	Size   int64 `json:"size"`
	Height int64 `json:"height"`
}

func NewState(path string) *State {
	var state State
	name := "state"
	db, err := db.NewGoLevelDB(name, path)
	if err != nil {
		log.Fatalf("failed to create persistent state at %s: %w", path, err)
	}
	stateBytes, err := db.Get(stateKey)
	if err != nil {
		log.Fatalf("failed to load state: %w", err)
	}
	if len(stateBytes) > 0 {
		err = json.Unmarshal(stateBytes, &state)
		if err != nil {
			log.Fatalf("failed to read current state: %w", err)
		}
	}
	state.db = db
	return &state
}

func (s *State) Save() error {
	stateBytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = s.db.Set(stateKey, stateBytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) Hash() []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(s.Size))
	return bytes
}

func (s *State) Close() {
	if err := s.db.Close(); err != nil {
		log.Printf("Closing state database: %v", err)
	}
}
