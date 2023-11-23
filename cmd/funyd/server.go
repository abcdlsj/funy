package main

import "github.com/abcdlsj/funy/internal/orchestrator"

func main() {
	orch := orchestrator.New()

	orch.Serve()
}
