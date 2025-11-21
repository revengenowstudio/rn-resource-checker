package main

import (
	"os"

	"rn-resource-checker/src/cli"
	"rn-resource-checker/src/log"
)

func main() {
	log.Init()
	cli.Run(os.Args)
}
