package orchestrator

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/abcdlsj/funy/internal/development"
	"github.com/abcdlsj/funy/internal/tarball"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Instance     development.Service `json:"-"`
	ProcessState ProcessState        `json:"process_state"`
	TarDir       string              `json:"tar_dir"`
	MainFile     string              `json:"main_file"`
	LDFlagX      map[string]string   `json:"ld_flag_x"`
}

type Orchestrator struct {
	mu       sync.Mutex
	Services []*Service
	Waitings []*Service
}

func New() *Orchestrator {
	return &Orchestrator{
		Services: make([]*Service, 0),
		Waitings: make([]*Service, 0),
	}
}

func (s *Orchestrator) Serve() {
	go s.Worker()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		ret := gin.H{}

		for _, service := range s.Services {
			ret[service.Instance.Name] = service
		}

		c.JSON(200, ret)
	})

	r.POST("/:service/create", func(c *gin.Context) {
		serviceName := c.Param("service")
		var req CreateRequest

		if err := c.BindJSON(&req); err != nil {
			c.String(400, err.Error())
			return
		}

		temp, err := os.MkdirTemp("", "funy*")
		if err != nil {
			c.String(500, err.Error())
			return
		}

		tarDir := filepath.Join(temp, serviceName)

		if err := os.Mkdir(tarDir, 0755); err != nil {
			c.String(500, err.Error())
			return
		}

		s.AddService(Service{
			Instance: development.Service{
				Name: serviceName,
			},
			ProcessState: Create,
			TarDir:       tarDir,
			MainFile:     req.MainFile,
			LDFlagX:      req.LDFlagX,
		})

		c.String(200, "create service success")
	})

	r.POST("/:service/deploy", func(c *gin.Context) {
		serviceName := c.Param("service")
		service := s.GetService(serviceName)

		if service == nil {
			c.String(404, "service not found")
			return
		}

		if err := tarball.Untar(service.TarDir, c.Request.Body); err != nil {
			c.String(500, err.Error())
			return
		}

		log.Infof("untar done: %s", service.TarDir)

		if service.ProcessState == Deployed {
			c.String(400, "service is already deployed")
			return
		}

		if service.ProcessState != Create {
			c.String(400, "service is deploying")
			return
		}

		service.ProcessState = Queued
		s.PushQueue(service)

		c.String(200, "will queue deploy service...")
	})

	r.Run(":8080")
}
