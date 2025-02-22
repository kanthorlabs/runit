package dockerx

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/kanthorlabs/runit/runtime/pythonx"
)

func TestExec(t *testing.T) {
	port := fmt.Sprintf("%d", rand.Intn(1000)+30000)
	vars := pythonx.NewDockerfileVars()
	vars.Ports = []string{port}

	err := Exec("../../examples/python/ip-checker.py", vars)
	if err != nil {
		t.Fatal(err)
	}
}
