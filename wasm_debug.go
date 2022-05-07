//go:build debug

package main

import (
	"net/http"
	"os"

	"github.com/soypat/gwasm"
)

func WASMHandler() (http.Handler, error) {
	wsm, err := gwasm.NewWASMHandler(appFolder, nil)
	if err != nil {
		return err
	}
	wsm.WASMReload = true
	wsm.SetOutput(os.Stdout)
	return wsm
}
