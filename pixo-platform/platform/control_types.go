package platform

import (
	"context"
	"encoding/json"
)

type ControlType struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GetControlTypesResponse struct {
	ControlTypes []*ControlType `json:"controls,omitempty"`
}

func (p *PlatformClient) GetControlTypes(ctx context.Context) ([]*ControlType, error) {
	query := `{ controls { id name } }`

	res, err := p.Client.ExecRaw(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var gqlRes GetControlTypesResponse
	if err = json.Unmarshal(res, &gqlRes); err != nil {
		return nil, err
	}

	return gqlRes.ControlTypes, nil
}
