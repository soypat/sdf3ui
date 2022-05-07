//go:build debug

package main

import (
	"net/http"
	"os"

	"github.com/soypat/gwasm"
)

func WASMHandler() (http.Handler, error) {
	wsm, err := gwasm.NewWASMHandler("app", nil)
	if err != nil {
		return nil, err
	}
	wsm.WASMReload = true
	wsm.SetOutput(os.Stdout)
	return wsm, nil
}
