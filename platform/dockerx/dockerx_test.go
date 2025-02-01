package dockerx

import "testing"

func TestExec(t *testing.T) {
	err := Exec("../../examples/python/ip-checker.py")
	if err != nil {
		t.Fatal(err)
	}
}
