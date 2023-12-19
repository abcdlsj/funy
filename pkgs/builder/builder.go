package builder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

type Builder struct {
	LDFlagX map[string]string
	WorkDir string
	Input   string
	Output  string
}

func (b *Builder) Build() error {
	args := []string{
		"build",
		"-buildmode=plugin",
		"-o",
		b.Output,
	}

	if len(b.LDFlagX) != 0 {
		args = append(args, "-ldflags")

		var ldflagx string
		for k, v := range b.LDFlagX {
			ldflagx += fmt.Sprintf("-X %s=%s ", k, v)
		}

		if len(ldflagx) != 0 {
			ldflagx = ldflagx[:len(ldflagx)-1]
		}

		args = append(args, ldflagx)
	}

	args = append(args, b.Input)

	log.Infof("builder args: %v", args)

	cmd := exec.Command("go", args...)
	cmd.Dir = b.WorkDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
