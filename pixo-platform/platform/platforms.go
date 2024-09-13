package platform

import (
	"context"
)

type Platform struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
}

type GetPlatformsResponse struct {
	Platforms []Platform `json:"platforms,omitempty"`
}

func (p *clientImpl) GetPlatforms(ctx context.Context) ([]Platform, error) {
	query := `query platforms { platforms { id name shortName } }`

	var res GetPlatformsResponse
	if err := p.Exec(ctx, query, &res, nil); err != nil {
		return nil, err
	}

	return res.Platforms, nil
}
