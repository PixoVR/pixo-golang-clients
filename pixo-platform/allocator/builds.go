package allocator

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/rs/zerolog/log"
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

	res, err := a.Get(path)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		log.Debug().Msgf("Received status code %d when getting build workflows: %s", res.StatusCode(), res.Body())
		return nil, fmt.Errorf("received status code %d when getting build workflows", res.StatusCode())
	}

	var workflowsResponse WorkflowsResponse
	if err = json.Unmarshal(res.Body(), &workflowsResponse); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal build workflows response")
		return nil, err
	}

	return workflowsResponse.Workflows, nil
}

// GetBuildWorkflowLogs returns a channel of logs for a build workflow or an error
func (a *Client) GetBuildWorkflowLogs(workflowName string) (chan *Log, error) {
	path := fmt.Sprintf("build/workflows/%s/logs", workflowName)

	req := a.FormatRequest()
	req.SetHeader("Accept", "application/octet-stream")

	req.SetDoNotParseResponse(true)

	res, err := req.Get(a.GetURLWithPath(path))
	if err != nil {
		return nil, err
	}

	logs := make(chan *Log, 100)

	go func() {
		for {
			line, err := bufio.NewReader(res.RawResponse.Body).ReadBytes('\n')
			if err != nil {
				close(logs)
				break
			}

			log.Debug().Msgf("Received line from stream: %s", line)

			data := strings.Trim(string(line), "\n")
			if data != "event:message" {
				var workflowLog Log

				workflowLog.Step = workflowName
				workflowLog.Lines = data
				//if err = json.Unmarshal(line, &workflowLog); err != nil {
				//	log.Debug().Err(err).Msg("Failed to unmarshal workflow log")
				//	continue
				//}

				logs <- &workflowLog
			}
		}
	}()

	return logs, nil
}
