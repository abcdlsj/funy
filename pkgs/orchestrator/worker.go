package orchestrator

import (
	"path/filepath"
	"time"

	"github.com/abcdlsj/funy/pkgs/builder"
	"github.com/abcdlsj/funy/pkgs/dyloader"
	"github.com/abcdlsj/funy/pkgs/share"
	"github.com/charmbracelet/log"
)

func (s *Orchestrator) ScheduleWorker() {
	deploy := func(ins *Instance) {
		dyPath := filepath.Join(share.DynamicLibDir, ins.InsName+".so")
		builder := builder.Builder{
			LDFlagX: ins.LDFlagX,
			WorkDir: ins.TarDir,
			Input:   ins.MainFile,
			Output:  dyPath,
		}

		if err := builder.Build(); err != nil {
			ins.ProcessState = share.Error
			log.Errorf("build error: %v", err)
			return
		}

		ins.ProcessState = share.Built

		if ins.AppType == "service" {
			runFn, shutdownFn, err := dyloader.LoadService(dyPath)

			if err != nil {
				ins.ProcessState = share.Error
				log.Errorf("load error: %v", err)
				return
			}

			ins.ProcessState = share.Loaded
			ins.Service.RunFn = runFn
			ins.Service.ShutdownFn = shutdownFn

			go func() {
				ins.ProcessState = share.Deployed

				if err := ins.Service.Run(); err != nil {
					ins.ProcessState = share.Error
					log.Errorf("run error: %v", err)
					return
				}
			}()
		} else {
			handlerFn, err := dyloader.LoadFunction(dyPath)
			if err != nil {
				ins.ProcessState = share.Error
				log.Errorf("load error: %v", err)
				return
			}

			ins.ProcessState = share.Loaded
			ins.Function.Handler = handlerFn

			s.fnRoutes.Store(ins.InsName, handlerFn)
			ins.ProcessState = share.Deployed
		}
	}

	for {
		service := s.Next()
		if service == nil {
			time.Sleep(5 * time.Second)
			continue
		}

		if service.ProcessState == share.Queued {
			go deploy(service)
		}
	}
}
