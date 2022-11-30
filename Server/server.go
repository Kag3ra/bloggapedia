package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/cors"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File\n")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprintln(w, "File too big")
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Fprintln(w, "Error receiving file form-data", err)
		return
	}

	defer file.Close()

	fmt.Printf("Upload File Name: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("/Users/mike/fun/fuzzy-wuzzy/Server/temp", "upload-*.jpg")
	if err != nil {
		fmt.Fprintln(w, err)
		return

	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(w, "Failed to read file", err)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		fmt.Fprintln(w, "Failed to write data to file", err)
		return
	}

	fmt.Fprintf(w, "File Uplooaded Successfully File\n")

	w.Header().Set("Location", "http://www.example.org/login")
	w.WriteHeader(200)
}

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./temp")
	if err != nil {
		fmt.Fprintln(w, "Failed to fetch posts", err)
		return
	}

	var posts = make([]string, len(files))
	for i, f := range files {
		posts[i] = f.Name()
	}
	body, err := json.Marshal(posts)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func setupRoutes() {
	mux := http.NewServeMux()
	port := ":8080"
	mux.HandleFunc("/upload", uploadFile)
	mux.HandleFunc("/posts", FetchPosts)
	handler := cors.Default().Handler(mux)
	fmt.Printf("I'm running at http://localhost%s\n", port)
	http.ListenAndServe(port, handler)
}

func main() {
	setupRoutes()
}
