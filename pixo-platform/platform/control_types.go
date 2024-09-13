package platform

import (
	"context"
)

type ControlType struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GetControlTypesResponse struct {
	ControlTypes []ControlType `json:"controls,omitempty"`
}

func (p *clientImpl) GetControlTypes(ctx context.Context) ([]ControlType, error) {
	query := `{ controls { id name } }`

	var res GetControlTypesResponse
	if err := p.Exec(ctx, query, &res, nil); err != nil {
		return nil, err
	}

	return res.ControlTypes, nil
}
