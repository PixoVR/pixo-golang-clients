package vertexai

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"os"
)

type Client struct {
	Response    Response
	projectID   string
	modelID     string
	accessToken string
	debug       bool
	httpClient  *resty.Client
}

type ClientConfig struct {
	ProjectID string
	ModelID   string
	Debug     bool
}

func NewClient(config ClientConfig) (*Client, error) {

	if config.ProjectID == "" {
		config.ProjectID = os.Getenv("GCP_PROJECT_ID")
		if config.ProjectID == "" {
			err := errors.New("GCP_PROJECT_ID is empty")
			log.Error().Err(err)
			return nil, err
		}
	}

	accessToken, err := getAccessToken()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &Client{
		debug:       config.Debug,
		projectID:   config.ProjectID,
		modelID:     config.ModelID,
		accessToken: accessToken,
		httpClient:  resty.New(),
	}, nil
}

func (c *Client) FormatVertexAPIURL(modelID string) string {
	return fmt.Sprintf("https://%s/v1/projects/%s/locations/us-central1/publishers/google/models/%s:predict", API_ENDPOINT, c.projectID, modelID)
}

type ChatRequest struct {
	Debug      bool
	Prompt     string
	ModelID    string
	Instances  []Instance
	Parameters *Parameters
}

func (c *Client) ChatResponse(req ChatRequest) (string, error) {

	if req.ModelID == "" {
		req.ModelID = c.modelID
		if req.ModelID == "" {
			req.ModelID = BISON_CODE_MODEL_ID
			c.debugMsgf("model id is empty, using default model id: %s", req.ModelID)
		} else {
			c.debugMsgf("model id is empty, using model id from config: %s", req.ModelID)
		}
	}

	if len(req.Instances) == 0 {
		if req.Prompt == "" {
			err := errors.New("instances and prompt cannot be empty")
			log.Error().Err(err)
			return "", err
		}

		c.debugMsg("instances is empty, using prompt as instance")
		req.Instances = []Instance{{Prefix: req.Prompt}}
	}

	if req.Parameters == nil {
		req.Parameters = &DEFAULT_PARAMETER

		c.debugMsg("parameters is empty, using default parameter")
	}

	url := c.FormatVertexAPIURL(req.ModelID)
	c.debugMsgf("client api url: %s", url)

	request := Request{
		Instances:  req.Instances,
		Parameters: *req.Parameters,
	}
	c.debugMsgf("client request: %+v", request)

	response, err := c.httpClient.R().
		SetAuthToken(c.accessToken).
		SetBody(request).
		Post(url)

	if err != nil {
		log.Err(err).Msg("could not get response from client api")
		return "", err
	}

	if err = json.Unmarshal(response.Body(), &c.Response); err != nil {
		log.Err(err).Msg("could not unmarshal response from client ai api")
		return "", err
	}

	if c.Response.Error.Code != 0 {
		msg := fmt.Sprintf("%d %s %s", c.Response.Error.Code, c.Response.Error.Status, c.Response.Error.Message)
		err = errors.New(msg)
		log.Err(err)
		return "", err
	}

	if len(c.Response.Predictions) == 0 {
		err = errors.New("no predictions returned")
		log.Err(err)
		return "", err
	}

	c.debugMsgf("vertexai response: %+v", c.Response)

	return c.Response.Predictions[0].Content, nil
}
