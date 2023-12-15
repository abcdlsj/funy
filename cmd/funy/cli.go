package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/abcdlsj/funy/pkgs/orchestrator"
	"github.com/abcdlsj/funy/pkgs/tarball"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var Version = ""

func main() {
	app := &cli.App{
		Name:  "funy",
		Usage: "development tool",
		Commands: []*cli.Command{
			{
				Name:  "app",
				Usage: "development app",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "server_address",
						Usage:    "server address",
						Value:    "127.0.0.1:8080",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "type",
						Usage:    "app type",
						Value:    "service",
						Required: false,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name: "create",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "service/function name",
								Required: true,
							},
							&cli.StringSliceFlag{
								Name:  "ld_flag_x",
								Usage: "set ld flag x",
							},
							&cli.StringFlag{
								Name:  "main_file",
								Usage: "set go file",
								Value: "main.go",
							},
						},
						Usage: "create app",
						Action: func(c *cli.Context) error {
							serverAddr := c.String("server_address")
							appType := c.String("type")
							name := c.String("name")
							mainFile := c.String("main_file")
							ldFlagX := c.StringSlice("ld_flag_x")

							return Create(serverAddr, name, mainFile, appType, ldFlagX)
						},
					},
					{
						Name:  "deploy",
						Usage: "deploy app",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "app name",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "dir",
								Usage:    "dir <path>",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							serverAddr := c.String("server_address")
							name := c.String("name")
							dir := c.String("dir")

							return Deploy(serverAddr, name, dir)
						},
					},
				},
			},
		},
		Action: func(*cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func authheader() string {
	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	pair := fmt.Sprintf("%s:%s", user, password)

	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(pair)))
}

func Create(svr, name, mainFile, appType string, ldFlagX []string) error {
	ldFlagXMap := make(map[string]string)

	for _, v := range ldFlagX {
		kv := strings.Split(v, "=")
		ldFlagXMap[kv[0]] = kv[1]
	}

	reqStru := orchestrator.CreateReq{
		MainFile: mainFile,
		LDFlagX:  ldFlagXMap,
		AppType:  appType,
	}

	body, err := json.Marshal(reqStru)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/orchestrator/%s/create", svr, name), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authheader())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func Deploy(svr string, name string, dir string) error {
	tar, err := tarball.TarDir(context.Background(), dir)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/orchestrator/%s/deploy", svr, name), tar)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-tar")
	req.Header.Set("Authorization", authheader())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
