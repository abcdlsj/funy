package orchestrator

import (
	"os"

	"github.com/charmbracelet/log"
)

const (
	DynamicLibDir = "/tmp/funydynamic/"
)

func init() {
	if err := os.MkdirAll(DynamicLibDir, 0755); err != nil {
		log.Fatalf("create dynamic lib dir error: %v", err)
	}
}

type ProcessState uint8

const (
	Create ProcessState = iota
	Queued
	Built
	Loaded
	Deployed

	Error
	Deleted
)

func (p ProcessState) String() string {
	switch p {
	case Create:
		return "CREATE"
	case Queued:
		return "QUEUED"
	case Built:
		return "BUILT"
	case Loaded:
		return "LOADED"
	case Deployed:
		return "DEPLOYED"
	case Error:
		return "ERROR"
	case Deleted:
		return "DELETED"
	}

	return "UNKNOWN"
}

type CreateRequest struct {
	MainFile string            `json:"main_file"`
	LDFlagX  map[string]string `json:"ld_flag_x"`
}
