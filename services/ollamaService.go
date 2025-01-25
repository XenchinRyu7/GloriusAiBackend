package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

var (
	activeModelName string
	modelMutex      sync.Mutex
	cacheMutex      sync.Mutex
	availableModels []map[string]string
)

type SendMessageRequest struct {
	Message string `json:"message"`
	Stream  bool   `json:"stream"`
}

func ActivateModel(modelName string) error {
	modelMutex.Lock()
	defer modelMutex.Unlock()

	if activeModelName == modelName {
		log.Printf("Model %s is already active.", modelName)
		return nil
	}

	ollamaPath := "C:\\Users\\user\\AppData\\Local\\Programs\\Ollama\\ollama.exe"
	cmd := exec.Command(ollamaPath, "run", modelName)

	log.Printf("Activating model: %s", modelName)
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to activate model: %s", err)
		return err
	}

	activeModelName = modelName
	log.Printf("Model %s activated successfully.", modelName)
	return nil
}

func GetAllModels() ([]map[string]string, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if availableModels != nil {
		log.Println("Returning cached model list.")
		return availableModels, nil
	}

	ollamaPath := "C:\\Users\\user\\AppData\\Local\\Programs\\Ollama\\ollama.exe"
	cmd := exec.Command(ollamaPath, "list")
	log.Printf("Fetching model list: %s %v", cmd.Path, cmd.Args)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to list models: %s", err)
		return nil, err
	}

	log.Println("Model list output:", string(output))
	availableModels, err = parseOllamaListOutput(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse model list: %w", err)
	}

	return availableModels, nil
}

func parseOllamaListOutput(output string) ([]map[string]string, error) {
	var models []map[string]string
	scanner := bufio.NewScanner(strings.NewReader(output))
	isHeader := true

	for scanner.Scan() {
		line := scanner.Text()
		if isHeader {
			isHeader = false
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		model := map[string]string{
			"name":     fields[0],
			"id":       fields[1],
			"size":     fields[2],
			"modified": strings.Join(fields[3:], " "),
		}
		models = append(models, model)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

type OllamaResponse struct {
	Model      string `json:"model"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason"`
}

func SendMessageToModel(jsonRequest []byte, w http.ResponseWriter) (string, error) {
	var req SendMessageRequest
	if err := json.Unmarshal(jsonRequest, &req); err != nil {
		return "", fmt.Errorf("failed to parse JSON request: %w", err)
	}

	if req.Stream {
		streamToOllamaAPI(activeModelName, req.Message, w)
		return "", nil
	}

	return sendToOllamaAPI(activeModelName, req.Message, false)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	response, err := SendMessageToModel(bodyBytes, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For non-streaming, return the response as JSON
	if !req.Stream {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"response": response})
	}
}

func sendToOllamaAPI(model, message string, stream bool) (string, error) {
	apiURL := "http://localhost:11434/api/generate"
	requestBody := map[string]interface{}{
		"model":  model,
		"prompt": message,
		"stream": stream,
	}

	reqBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request to Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error: %s", string(body))
	}

	if stream {
		return "", errors.New("streaming response not supported in this function")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var response OllamaResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Done {
		return "", fmt.Errorf("ollama API returned unfinished response: %s", response.DoneReason)
	}

	return response.Response, nil
}

func streamToOllamaAPI(model, message string, w http.ResponseWriter) {
	apiURL := "http://localhost:11434/api/generate"
	requestBody := map[string]interface{}{
		"model":  model,
		"prompt": message,
		"stream": true,
	}

	reqBytes, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal request body: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send HTTP request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Ollama API error: %s", string(body)), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(w, "data: %s\n\n", line)
		w.(http.Flusher).Flush()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error streaming response: %v", err)
	}
}
