package allocator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/argo"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

type WorkflowsResponse struct {
	Message   string
	Workflows []Workflow
}

type Workflow struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (a *AllocatorClient) GetBuildWorkflows() ([]Workflow, error) {
	path := "build/workflows"

	res, err := a.Get(path)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf("received status code %d when getting build workflows", res.StatusCode())
		log.Debug().Err(err).Msg("Failed to get build workflows")
		return nil, err
	}

	var workflowsResponse WorkflowsResponse
	if err = json.Unmarshal(res.Body(), &workflowsResponse); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal get build workflows response")
		return nil, err
	}

	return workflowsResponse.Workflows, nil
}

func (a *AllocatorClient) GetBuildWorkflowLogs(workflowName string) (chan *argo.Log, error) {
	path := fmt.Sprintf("build/workflows/%s/logs", workflowName)

	req := a.FormatRequest()
	req.SetHeader("Accept", "application/octet-stream")

	req.SetDoNotParseResponse(true)

	res, err := req.Get(a.GetURLWithPath(path))
	if err != nil {
		return nil, err
	}

	logs := make(chan *argo.Log, 100)

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
				var workflowLog argo.Log

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
