package builder_test

import (
	"testing"

	"github.com/abcdlsj/funy/internal/builder"
)

func TestBuiler(t *testing.T) {
	bu := builder.Builder{
		LDFlagX: map[string]string{
			"main.Version": "1.0.0",
		},
		WorkDir: "../../cmd/funy/",
		Input:   "cli.go",
		Output:  "/tmp/funy_cli_plugin.so",
	}

	if err := bu.Build(); err != nil {
		t.Fatal(err)
	}

	t.Logf("Output: %s", bu.Output)
}
