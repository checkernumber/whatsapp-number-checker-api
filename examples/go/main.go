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

var apiKey = mustEnv("CHECKNUMBER_API_KEY")

type Task struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	Total     int    `json:"total"`
	Success   int    `json:"success"`
	Failure   int    `json:"failure"`
	ResultURL string `json:"result_url"`
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Fprintf(os.Stderr, "Set the %s environment variable\n", k)
		os.Exit(1)
	}
	return v
}

func download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("download failed: http %d", resp.StatusCode)
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
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
	if task.ResultURL == "" {
		fmt.Println("task exported but no result_url")
		return
	}
	if err := download(task.ResultURL, "results.zip"); err != nil {
		panic(err)
	}
	fmt.Println("saved to: results.zip")
}
