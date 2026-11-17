package platform

import "context"

type License struct {
	ID   int    `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func (p *clientImpl) GetLicensesForUser(ctx context.Context, userID int) ([]License, error) {
	query := `query user($id: ID!){ user(id: $id) { licenses { id type name } } }`

	variables := map[string]interface{}{
		"id": userID,
	}

	var userResponse GetUserResponse
	if err := p.Exec(ctx, query, &userResponse, variables); err != nil {
		return nil, err
	}

	return userResponse.User.Licenses, nil
}
