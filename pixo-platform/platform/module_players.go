package platform

import (
	"context"
	"time"
)

type ModulePlayer struct {
	ID             int       `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	DistributorID  int       `json:"distributorId,omitempty"`
	LaunchProtocol string    `json:"launchProtocol,omitempty"`
	Distributor    Org       `json:"distributor,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`

	Versions []ModulePlayerVersion `json:"versions,omitempty"`
}

type ModulePlayerParams struct {
	UserID *int `json:"userId,omitempty"`
}

type GetModulePlayersResponse struct {
	ModulePlayers []ModulePlayer `json:"modulePlayers"`
}

func (p *clientImpl) GetModulePlayers(ctx context.Context, params ...*ModulePlayerParams) ([]ModulePlayer, error) {
	query := `query modulePlayers($params: ModulePlayerParamsInput) { modulePlayers(params: $params) { id name description distributorId launchProtocol distributor { id name } createdAt updatedAt } }`

	var p0 *ModulePlayerParams
	if len(params) > 0 {
		p0 = params[0]
	}
	variables := map[string]interface{}{
		"params": p0,
	}

	var res GetModulePlayersResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.ModulePlayers, nil
}
