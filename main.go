package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(Data)
	systray.SetTitle("")
	systray.SetTooltip("Minio Server")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	cmd := runMinio()

	go func() {
		<-mQuit.ClickedCh
		cmd.Process.Kill()
		os.Exit(0)
	}()
}

func runMinio() *exec.Cmd {
	DownloadFile("/tmp/minio-server", "https://dl.min.io/server/minio/release/darwin-amd64/minio")
	os.Chmod("/tmp/minio-server", 0777)

	cmd := exec.Command("/tmp/minio-server", "server", "--address", ":19000", "/tmp/.minio-server-data")
	env := append(os.Environ(), []string{
		"MINIO_ACCESS_KEY=asdfasdf",
		"MINIO_SECRET_KEY=qwerqwer",
	}...)
	cmd.Env = env
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Start()
	fmt.Println("Started")
	return cmd
}

func onExit() {
	// clean up here
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
