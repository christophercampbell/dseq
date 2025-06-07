package app

import (
	"github.com/cometbft/cometbft/abci/types"
)

// StateManager defines the interface for managing application state
type StateManager interface {
	Save() error
	Hash() []byte
	Close() error
}

// Sequencer defines the interface for the sequencer application
type Sequencer interface {
	types.Application
	GetID() string
	GetAddress() string
	GetState() StateManager
}

// StreamManager defines the interface for managing data streams
type StreamManager interface {
	Start() error
	Stop() error
	ProcessEntry(entry interface{}) error
}

// Logger defines the interface for logging
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	With(keyvals ...interface{}) Logger
}

// ConfigManager defines the interface for managing configuration
type ConfigManager interface {
	Load() error
	Validate() error
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}
