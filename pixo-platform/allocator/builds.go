package allocator

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WorkflowsResponse represents the response from the build/workflows endpoint
type WorkflowsResponse struct {
	Message   string     `json:"message"`
	Workflows []Workflow `json:"workflows"`
}

// Workflow represents a build workflow
type Workflow struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// Log represents a log line from a build workflow
type Log struct {
	Step  string `json:"step"`
	Lines string `json:"lines"`
}

// GetBuildWorkflows returns a list of build workflows
func (a *Client) GetBuildWorkflows() ([]Workflow, error) {
	path := "build/workflows"

	res, err := a.Get(context.TODO(), path)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d when getting build workflows", res.StatusCode)
	}

	resBody, _ := io.ReadAll(res.Body)

	var workflowsResponse WorkflowsResponse
	if err = json.Unmarshal(resBody, &workflowsResponse); err != nil {
		return nil, err
	}

	return workflowsResponse.Workflows, nil
}

// GetBuildWorkflowLogs returns a channel of logs for a build workflow or an error
func (a *Client) GetBuildWorkflowLogs(workflowName string) (chan *Log, error) {
	path := fmt.Sprintf("build/workflows/%s/logs", workflowName)

	req, err := a.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")

	res, err := a.Client().Do(req)
	if err != nil {
		return nil, err
	}

	logs := make(chan *Log, 100)

	go func() {
		for {
			line, err := bufio.NewReader(res.Body).ReadBytes('\n')
			if err != nil {
				close(logs)
				break
			}

			data := strings.Trim(string(line), "\n")
			if data != "event:message" {
				var workflowLog Log

				workflowLog.Step = workflowName
				workflowLog.Lines = data
				//if err = json.Unmarshal(line, &workflowLog); err != nil {
				//	continue
				//}

				logs <- &workflowLog
			}
		}
	}()

	return logs, nil
}
