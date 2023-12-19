package orchestrator

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/abcdlsj/funy/pkgs/share"
	"github.com/abcdlsj/funy/pkgs/tarball"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type Orchestrator struct {
	Actives  []*Instance
	Waitings []*Instance

	fnRoutes sync.Map

	mu sync.Mutex
}

func New() *Orchestrator {
	return &Orchestrator{
		Actives:  make([]*Instance, 0),
		Waitings: make([]*Instance, 0),
	}
}

func (s *Orchestrator) Serve() {
	go s.ScheduleWorker()

	r := gin.Default()

	r.Use(gin.BasicAuth(gin.Accounts{
		os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
	}))

	r.GET("/secret", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"secret": "The secret ingredient to the BBQ sauce is stiring it in an old whiskey barrel.",
		})
	})

	org := r.Group("/orchestrator")

	org.GET("/", func(c *gin.Context) {
		ret := gin.H{}

		for _, v := range s.GetAll() {
			ret[v.InsName] = v
		}

		c.JSON(200, ret)
	})

	org.POST("/:instance/create", func(c *gin.Context) {
		insName := c.Param("instance")
		var req CreateReq

		if err := c.BindJSON(&req); err != nil {
			c.String(400, err.Error())
			return
		}

		temp, err := os.MkdirTemp("", "funy*")
		if err != nil {
			c.String(500, err.Error())
			return
		}

		tarDir := filepath.Join(temp, insName)

		if err := os.Mkdir(tarDir, 0755); err != nil {
			c.String(500, err.Error())
			return
		}

		s.Add(Instance{
			AppType:      req.AppType,
			InsName:      insName,
			ProcessState: share.Create,
			TarDir:       tarDir,
			MainFile:     req.MainFile,
			LDFlagX:      req.LDFlagX,
		})

		c.String(200, "create service success")
	})

	org.POST("/:instance/deploy", func(c *gin.Context) {
		insName := c.Param("instance")
		instance := s.GetByName(insName)

		if instance == nil {
			c.String(404, "service not found")
			return
		}

		if err := tarball.Untar(instance.TarDir, c.Request.Body); err != nil {
			c.String(500, err.Error())
			return
		}

		log.Infof("untar done: %s", instance.TarDir)

		if instance.ProcessState == share.Deployed {
			c.String(400, "service is already deployed")
			return
		}

		if instance.ProcessState != share.Create {
			c.String(400, "service is deploying")
			return
		}

		instance.ProcessState = share.Queued
		s.Push(instance)

		c.String(200, "will queue deploy service...")
	})

	fng := r.Group("/func")

	fng.GET("/", func(c *gin.Context) {
		routes := make([]string, 0)
		s.fnRoutes.Range(func(k, v interface{}) bool {
			routes = append(routes, k.(string))
			return true
		})

		c.JSON(200, routes)
	})

	fng.GET("/:function/*action", func(c *gin.Context) {
		fnName := c.Param("function")
		fn, ok := s.fnRoutes.Load(fnName)
		if !ok {
			c.String(404, "function not found")
			return
		}

		gin.WrapF(fn.(func(http.ResponseWriter, *http.Request)))(c)
	})

	// fng.POST("/:function/*action", func(c *gin.Context) {
	// 	fnName := c.Param("function")
	// 	fn, ok := s.fnRoutes.Load(fnName)
	// 	if !ok {
	// 		c.String(404, "function not found")
	// 		return
	// 	}

	// 	fn.(func(http.ResponseWriter, *http.Request))(c.Writer, c.Request)
	// })

	r.Run(":" + os.Getenv("PORT"))
}
