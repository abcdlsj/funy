package share

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
