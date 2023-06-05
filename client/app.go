package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	ee *exec.ExitError
	pe *os.PathError
)

const (
	downloadScript string = "http://localhost:8081/download/"
	uploadScript   string = "http://localhost:8081/upload"
)

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func makeLinuxScreenshot(tmpFilename *os.File) *os.File {
	collection := map[string][]string{
		"gnome-screenshot": {"-a", "-f", tmpFilename.Name()},
		"mv":               {"$(xfce4-screenshooter -r -o ls)", tmpFilename.Name()},
		"spectacle":        {"-b", "-n", "-r -o", tmpFilename.Name()},
		"scrot":            {"-s", tmpFilename.Name()},
		"import":           {tmpFilename.Name()},
	}
	for command, params := range collection {
		cmd := exec.Command(command, params...)
		_, err := cmd.CombinedOutput()
		if errors.As(err, &ee) {
			log.Println("exit code error:", ee.ExitCode()) // run, !=0 exit code
			continue

		} else if errors.As(err, &pe) {
			log.Printf("os.PathError: %v", pe) // "no such file ...", "permission denied" etc.
			continue

		} else if err != nil {
			log.Printf("general error: %v", err) // something errors!
			continue

		} else {
			log.Println("success!") // run ==0 exit code
			break
		}
	}
	return tmpFilename
}

// makeDarwinScreenshot macOS screenshot
func makeDarwinScreenshot(tmpFilename *os.File) *os.File {
	collection := map[string][]string{
		"screencapture": {"-i", tmpFilename.Name()},
	}
	for command, params := range collection {
		cmd := exec.Command(command, params...)
		_, err := cmd.CombinedOutput()
		if errors.As(err, &ee) {
			log.Println("exit code error:", ee.ExitCode()) // run, !=0 exit code
			continue

		} else if errors.As(err, &pe) {
			log.Printf("os.PathError: %v", pe) // "no such file ...", "permission denied" etc.
			continue

		} else if err != nil {
			log.Printf("general error: %v", err) // something errors!
			continue

		} else {
			log.Println("success!") // run ==0 exit code
			break
		}
	}
	return tmpFilename
}

// makeWindowsScreenshot '/clip' requires at least Win10 1703
func makeWindowsScreenshot(tmpFilename *os.File) *os.File {
	collection := map[string][]string{
		"snippingtool": {"/clip", tmpFilename.Name()},
	}
	for command, params := range collection {
		cmd := exec.Command(command, params...)
		_, err := cmd.CombinedOutput()
		if errors.As(err, &ee) {
			log.Println("exit code error:", ee.ExitCode()) // run, !=0 exit code
			continue

		} else if errors.As(err, &pe) {
			log.Printf("os.PathError: %v", pe) // "no such file ...", "permission denied" etc.
			continue

		} else if err != nil {
			log.Printf("general error: %v", err) // something errors!
			continue

		} else {
			log.Println("success!") // run ==0 exit code
			break
		}
	}
	return tmpFilename
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

type fileName struct {
	Name string `json:"name"`
}

func uploadFile(file *os.File) string {

	request, err := newfileUploadRequest(uploadScript,
		nil,
		"fileupload",
		file.Name())
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	data := fileName{}
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	url := downloadScript + data.Name
	return url

}

func main() {
	osName := runtime.GOOS
	tmpFile, err := os.CreateTemp("", "screenshot*.png")
	checkErr(err)
	log.Printf("Temp file created:", tmpFile.Name())
	defer os.Remove(tmpFile.Name())
	log.Printf(tmpFile.Name())
	switch osName {
	case "windows":
		screenshot := makeWindowsScreenshot(tmpFile)
		fileName := uploadFile(screenshot)
		openBrowser(fileName)
	case "darwin":
		screenshot := makeDarwinScreenshot(tmpFile)
		fileName := uploadFile(screenshot)
		openBrowser(fileName)
	case "linux":
		screenshot := makeLinuxScreenshot(tmpFile)
		fileName := uploadFile(screenshot)
		openBrowser(fileName)
	default:
		fmt.Printf("%s.\n", osName)
	}
}
