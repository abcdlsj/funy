package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var Srv = &http.Server{
	Addr: ":" + Port,
}

var Port string

var Count = int32(0)

func Run() error {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		atomic.AddInt32(&Count, 1)

		c.JSON(200, gin.H{
			"ping count": fmt.Sprint(atomic.LoadInt32(&Count)),
		})
	})

	Srv.Handler = r

	if err := Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("listen: %s\n", err)
		return err
	}

	return nil
}

func Shutdown() error {
	Srv.Shutdown(context.Background())
	log.Println("Shutdown Server ...")

	return nil
}
