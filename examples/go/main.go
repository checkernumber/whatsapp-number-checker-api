// WhatsApp Number Checker API — bulk verification example (Go).
// Workflow: submit a file of E.164 numbers -> poll status -> download results.
// Docs: https://docs.checknumber.ai/whatsapp-bulk-checker
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	baseURL  = "https://api.checknumber.ai"
	taskType = "ws" // ws | ws_active | ws_avatar
)

var apiKey = getenv("CHECKNUMBER_API_KEY", "YOUR_API_KEY")

type Task struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	Total     int    `json:"total"`
	Success   int    `json:"success"`
	Failure   int    `json:"failure"`
	ResultURL string `json:"result_url"`
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func submitTask(path string) (*Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile("file", path)
	io.Copy(part, file)
	w.WriteField("task_type", taskType)
	w.Close()

	req, _ := http.NewRequest("POST", baseURL+"/v1/tasks", &body)
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return doTask(req)
}

func getTask(taskID string) (*Task, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField("task_id", taskID)
	w.Close()

	req, _ := http.NewRequest("POST", baseURL+"/v1/gettasks", &body)
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return doTask(req)
}

func doTask(req *http.Request) (*Task, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}
	var t Task
	return &t, json.NewDecoder(resp.Body).Decode(&t)
}

func main() {
	task, err := submitTask("numbers.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("task_id:", task.TaskID)

	for {
		task, err = getTask(task.TaskID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("status=%s success=%d/%d\n", task.Status, task.Success, task.Total)
		if task.Status == "exported" {
			break
		}
		if task.Status == "failed" {
			panic("task failed")
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("result_url:", task.ResultURL)
}
