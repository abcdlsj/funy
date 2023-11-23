package orchestrator

import (
	"path/filepath"
	"time"

	"github.com/abcdlsj/funy/internal/builder"
	"github.com/abcdlsj/funy/internal/dyloader"
	"github.com/charmbracelet/log"
)

func (s *Orchestrator) Worker() {
	deploy := func(service *Service) {
		dyPath := filepath.Join(DynamicLibDir, service.Instance.Name+".so")
		builder := builder.Builder{
			LDFlagX: service.LDFlagX,
			WorkDir: service.TarDir,
			Input:   service.MainFile,
			Output:  dyPath,
		}

		if err := builder.Build(); err != nil {
			service.ProcessState = Error
			log.Errorf("build error: %v", err)
			return
		}

		service.ProcessState = Built

		runFn, shutdownFn, err := dyloader.Load(dyPath)

		if err != nil {
			service.ProcessState = Error
			log.Errorf("load error: %v", err)
			return
		}

		service.ProcessState = Loaded
		service.Instance.RunFn = runFn
		service.Instance.ShutdownFn = shutdownFn

		go func() {
			service.ProcessState = Deployed

			if err := service.Instance.Run(); err != nil {
				service.ProcessState = Error
				log.Errorf("run error: %v", err)
				return
			}
		}()
	}

	for {
		service := s.NextQueue()
		if service == nil {
			time.Sleep(5 * time.Second)
			continue
		}

		if service.ProcessState == Queued {
			go deploy(service)
		}
	}
}
