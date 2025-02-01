package main

import (
	_ "embed"
	"log"

	"github.com/kanthorlabs/runit/platform/dockerx"
)

func main() {
	err := dockerx.Exec("../../examples/python/ip-checker.py")
	if err != nil {
		log.Fatal(err)
	}
}
