package legacy

import (
	"encoding/json"
	"errors"
)

// GetOrgs returns a list of webhooks.
func (p *Client) GetOrgs() ([]Org, error) {
	url := p.GetURLWithPath("api/orgs")

	res, err := p.FormatRequest().
		Get(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, errors.New(string(res.Body()))
	}

	var orgsRes struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Orgs    []Org  `json:"orgs"`
	}

	if err = json.Unmarshal(res.Body(), &orgsRes); err != nil {
		return nil, err
	}

	return orgsRes.Orgs, nil
}
