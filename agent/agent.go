package agent

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Run() {
	t0 := time.Now().UnixMilli()
	if len(os.Args) < 3 || os.Args[1] == "" && os.Args[2] == "" {
		fmt.Println("Please provide respectively the log folder and the openlog host as argument")
		os.Exit(1)
	}

	walkThroughFiles()

	t1 := time.Now().UnixMilli()
	fmt.Printf("The program took %v milliseconds to execute\n", t1-t0)
	wg.Wait()
}

func walkThroughFiles() {
	folder := os.Args[1]

	filepath.Walk(folder, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
		}
		if !file.IsDir() {
			go sendFileToOpenlog(path)
		}
		return nil
	})
}

func sendFileToOpenlog(path string) {
	wg.Add(1)

	request := createMultipartRequest(path)
	client := http.Client{}
	_, err := client.Do(request)

	if err != nil {
		fmt.Printf("Error sending file %s\n", path)
	}
	wg.Done()
}

func createMultipartRequest(filePath string) *http.Request {
	openlogHost := os.Args[2]

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s with error: %v\n", filePath, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("logfile", filePath)
	if err != nil {
		log.Printf("Error creating form file for %s with error: %v\n", filePath, err)
	}

	io.Copy(part, file)
	writer.Close()

	url := fmt.Sprintf("%s/openlog/api/v1/logcore/csv", openlogHost)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("Error creating http request for file %s with error: %v\n", filePath, err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request
}
