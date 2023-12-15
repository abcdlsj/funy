package main

import "github.com/abcdlsj/funy/pkgs/orchestrator"

func main() {
	orch := orchestrator.New()

	orch.Serve()
}
