package main

import (
	"github.com/MAHMETT/dockkit/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
