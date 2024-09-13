package platform

import (
	"context"
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

func (p *clientImpl) GetRoles(ctx context.Context) ([]Role, error) {
	query := `{ roles { id name permissions } }`

	var res GetRolesResponse
	if err := p.Exec(ctx, query, &res, nil); err != nil {
		return nil, err
	}

	return res.Roles, nil
}
