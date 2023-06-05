package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func createFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

type FileName struct {
	Name string `json:"name"`
}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	file, handler, err := r.FormFile("fileupload")
	if err != nil {
		return
	}
	defer file.Close()
	fileName := handler.Filename
	f, err := createFile("./save_files/" + fileName)
	if err != nil {
		return
	}
	if _, err = io.Copy(f, file); err != nil {
		return
	}
	defer f.Close()

	filename := FileName{
		Name: fileName,
	}
	json.NewEncoder(w).Encode(filename)
}

func Download(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	fp := path.Join("/save_files/", filename)
	http.ServeFile(w, r, fp)
}

func main() {
	router := mux.NewRouter()
	router.
		Path("/upload").
		Methods("POST").
		HandlerFunc(FileUpload)
	router.HandleFunc("/download/{filename}", Download).Name("download")
	fmt.Println("Starting")
	log.Fatal(http.ListenAndServe(":8081", router))
}
