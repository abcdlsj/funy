package orchestrator

import (
	"github.com/abcdlsj/funy/pkgs/development"
	"github.com/abcdlsj/funy/pkgs/share"
)

type CreateReq struct {
	MainFile string            `json:"main_file"`
	LDFlagX  map[string]string `json:"ld_flag_x"`
	AppType  string            `json:"app_type"`
}

type Instance struct {
	AppType      string               `json:"app_type"`
	InsName      string               `json:"instance_name"`
	Service      development.Service  `json:"-"`
	Function     development.Function `json:"-"`
	ProcessState share.ProcessState   `json:"process_state"`
	TarDir       string               `json:"tar_dir"`
	MainFile     string               `json:"main_file"`
	LDFlagX      map[string]string    `json:"ld_flag_x"`
}
