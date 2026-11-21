package allocator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type gameServersResponse struct {
	Results []GameServer `json:"results"`
}

func (a *Client) GetGameServers(ctx context.Context) ([]GameServer, error) {
	res, err := a.Get(ctx, "gameservers")
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d when getting game servers", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response gameServersResponse
	if err = json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}
