package main

import (
	"os"

	"github.com/kubeshop/testkube-executor-soapui/pkg/runner"
	"github.com/kubeshop/testkube/pkg/executor/agent"
)

func main() {
	agent.Run(runner.NewRunner(), os.Args)
}
