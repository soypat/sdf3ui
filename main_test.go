package main_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestVetApp(t *testing.T) {
	cmd := exec.Command("go", "vet", "./app")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(output), err)
	}
}
