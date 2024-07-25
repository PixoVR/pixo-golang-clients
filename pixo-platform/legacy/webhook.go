package legacy

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CreateWebhook creates a webhook with the given description and endpoint.
func (p *LegacyAPIClient) CreateWebhook(input Webhook) error {
	url := p.GetURLWithPath("api/webhook")

	res, err := p.FormatRequest().
		SetBody(input).
		Post(url)
	if err != nil {
		return err
	}

	if res.IsError() {
		return errors.New(string(res.Body()))
	}

	return nil
}

// GetWebhooks returns a list of webhooks.
func (p *LegacyAPIClient) GetWebhooks(orgID int) ([]Webhook, error) {
	url := p.GetURLWithPath(fmt.Sprintf("api/webhooks/org/%d", orgID))

	res, err := p.FormatRequest().
		Get(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, errors.New(string(res.Body()))
	}

	var webhooksRes struct {
		Error    bool      `json:"error"`
		Message  string    `json:"message"`
		Webhooks []Webhook `json:"webhooks"`
	}

	if err = json.Unmarshal(res.Body(), &webhooksRes); err != nil {
		return nil, err
	}

	return webhooksRes.Webhooks, nil
}

// DeleteWebhook deletes a webhook with the given ID.
func (p *LegacyAPIClient) DeleteWebhook(id int) error {
	url := p.GetURLWithPath(fmt.Sprintf("api/webhook/%d", id))

	res, err := p.FormatRequest().
		Delete(url)
	if err != nil {
		return err
	}

	if res.IsError() {
		return errors.New(string(res.Body()))
	}

	return nil
}
