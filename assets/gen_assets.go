package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const dirPerm = 0777

func main() {
	for _, dir := range []string{"js", "wasm"} {
		err := os.MkdirAll("assets/"+dir, dirPerm)
		if err != nil && !os.IsExist(err) {
			log.Fatal("creating directory: ", err)
		}
	}

	for _, w := range work {
		w.download()
	}
	const mainWasmPath = "assets/wasm/main.wasm"
	err := generateWASM(mainWasmPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("all files downloaded succesfully or already downloaded.")
	log.Println("program finished succesfully")
}

var work = []resource{
	{
		local:     "assets/js/trackball_controls.js",
		remoteURL: "https://raw.githubusercontent.com/soypat/three/main/examples/earth/assets/trackball_controls.js",
	},
}

type resource struct {
	remoteURL string
	local     string
}

func (r resource) download() {
	if _, err := os.Stat(r.local); err == nil {
		log.Printf("skipping: %s already exists\n", r.local)
		return
	}
	log.Println("download: ", r.local)
	err := downloadFile(r.remoteURL, r.local)
	if err != nil {
		log.Fatal("during file download: ", err)
	}
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func generateWASM(dst string) error {
	if _, err := os.Stat(dst); err == nil {
		log.Println(dst + " already generated. Delete file first to regenerate.")
		return nil
	}
	log.Println("generating WASM")
	// Build wasm file
	cmd := exec.Command("go", "build", "-o=main.wasm", ".")
	cmd.Dir = "app"
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		log.Print(string(out))
	}
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(cmd.Dir, "main.wasm"), dst)
	if err != nil {
		return err
	}
	return nil
}
