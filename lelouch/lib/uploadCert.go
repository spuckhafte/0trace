package lib

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// UploadCert uploads the wipe_log.json file to 0x0.st and returns the URL
func UploadCert(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	fileWriter, err := writer.CreateFormFile("file", "wipe_log.json")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the file content to the form
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://0x0.st", &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	// The response should contain the URL
	url := strings.TrimSpace(string(responseBody))
	return url, nil
}
