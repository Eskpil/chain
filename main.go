package main

import (
	"chain/cmd"
	"chain/logger"
)

func main() {
	logger.Init()
	cmd.Execute()
}
