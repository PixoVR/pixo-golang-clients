package platform

import (
	"context"
	"encoding/json"
	"time"
)

type Role struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Permissions string `json:"permissions,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type GetRolesResponse struct {
	Roles []Role `json:"roles"`
}

func (p *PlatformClient) GetRoles(ctx context.Context) ([]Role, error) {
	query := `{ roles { id name permissions } }`

	res, err := p.Client.ExecRaw(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var rolesResponse GetRolesResponse
	if err = json.Unmarshal(res, &rolesResponse); err != nil {
		return nil, err

	}

	return rolesResponse.Roles, nil
}
