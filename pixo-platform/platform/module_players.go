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
}

type GetModulePlayersResponse struct {
	ModulePlayers []ModulePlayer `json:"modulePlayers"`
}

func (p *clientImpl) GetModulePlayers(ctx context.Context) ([]ModulePlayer, error) {
	query := `query modulePlayers { modulePlayers { id name description distributorId launchProtocol distributor { id name } createdAt updatedAt } }`

	var res GetModulePlayersResponse
	if err := p.Exec(ctx, query, &res, nil); err != nil {
		return nil, err
	}

	return res.ModulePlayers, nil
}
