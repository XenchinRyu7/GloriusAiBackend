package services

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
)

var activeModelProcess *exec.Cmd
var modelStdin io.WriteCloser
var modelMutex sync.Mutex

func ActivateModel(modelName string) error {
	modelMutex.Lock()
	defer modelMutex.Unlock()

	// If an active process exists, terminate it
	if activeModelProcess != nil {
		activeModelProcess.Process.Kill()
	}

	// Start a new process for the specified model
	cmd := exec.Command("ollama", "run", modelName)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.New("failed to initialize stdin for the model: " + err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New("failed to initialize stdout for the model: " + err.Error())
	}

	if err := cmd.Start(); err != nil {
		return errors.New("failed to start model process: " + err.Error())
	}

	// Read output from the model in the background (optional)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Println(scanner.Text()) // Logging output for debugging
		}
	}()

	activeModelProcess = cmd
	modelStdin = stdin
	return nil
}

// SendMessageToModel sends a message to the active model and retrieves the response
func SendMessageToModel(message string) (string, error) {
	modelMutex.Lock()
	defer modelMutex.Unlock()

	if activeModelProcess == nil || modelStdin == nil {
		return "", errors.New("no active model process")
	}

	// Write message to the model's stdin
	_, err := modelStdin.Write([]byte(message + "\n"))
	if err != nil {
		return "", errors.New("failed to send message to model: " + err.Error())
	}

	// Read response from the model's stdout
	stdout, err := activeModelProcess.StdoutPipe()
	if err != nil {
		return "", errors.New("failed to read response from model: " + err.Error())
	}

	scanner := bufio.NewScanner(stdout)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", errors.New("no response from model")
}

// GetAllModels retrieves a list of all available models
func GetAllModels() ([]map[string]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("failed to execute ollama list command: " + err.Error())
	}

	return parseOllamaListOutput(string(output))
}

// parseOllamaListOutput parses the output of the `ollama list` command into a structured format
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
