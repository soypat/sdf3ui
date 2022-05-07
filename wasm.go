//go:build !debug

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/soypat/gwasm"
)

func WASMHandler() (http.Handler, error) {
	const compiler = "go"
	fp, err := assetsFS.Open("assets/wasm/main.wasm")
	if err != nil {
		return nil, fmt.Errorf("WASM application not found in executable: %w", err)
	}
	wasmApp, err := io.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("WASM application fail to read: %w", err)
	}
	// Load the Go wasmexec script.
	out, err := exec.Command(compiler, "env", "GOROOT").Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, string(out))
	}
	f := filepath.Join(strings.TrimSpace(string(out)), "misc", "wasm", "wasm_exec.js")
	fp, err = os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("WASM execution script open file error: %w", err)
	}
	wasmexec, err := io.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("reading WASM execution script: %w", err)
	}
	wsm := gwasm.WASMHandler{
		IndexHTML:       indexHTML,
		WASMReload:      true,
		Compiler:        compiler,
		WASMApplication: wasmApp,
		WASMExecContent: wasmexec,
	}
	wsm.SetOutput(os.Stdout)
	return &wsm, nil
}

const indexHTML = `<!DOCTYPE html>
<!-- Polyfill for the old Edge browser -->
<script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script>
<script src="wasm_exec.js"></script>
<script>
(async () => {
  const resp = await fetch('main.wasm');
  if (!resp.ok) {
    const pre = document.createElement('pre');
    pre.innerText = await resp.text();
    document.body.appendChild(pre);
  } else {
    const src = await resp.arrayBuffer();
    const go = new Go();
    const result = await WebAssembly.instantiate(src, go.importObject);
    go.argv = [];
    go.run(result.instance);
  }
  const reload = await fetch('_wait');
  // The server sends a response for '_wait' when a request is sent to '_notify'.
  if (reload.ok) {
    location.reload();
  }
})();
</script>
`
