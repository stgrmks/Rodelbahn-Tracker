package main

import "github.com/stgrmks/Rodelbahn-Tracker/cmd"

var (
	VERSION = "0.1.0"
	BUILD   = "0.1.0"
)

func main() {
	cmd.Execute(VERSION, BUILD)
}
