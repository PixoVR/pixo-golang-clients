package platform

import (
	"context"
	"errors"
)

type Webhook struct {
	ID            int      `json:"id,omitempty"`
	URL           string   `json:"url,omitempty"`
	EventTypes    []string `json:"eventTypes,omitempty"`
	Description   string   `json:"description,omitempty"`
	Token         string   `json:"token,omitempty"`
	GenerateToken *bool    `json:"generateToken,omitempty"`
	OrgID         int      `json:"orgId,omitempty"`
	Org           *Org     `json:"org,omitempty"`
}

type WebhookParams struct {
	OrgID int `json:"orgId"`
}

type GetWebhooksResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

type GetWebhookResponse struct {
	Webhook Webhook `json:"webhook"`
}

type CreateWebhookResponse struct {
	Webhook Webhook `json:"createWebhook"`
}

type UpdateWebhookResponse struct {
	Webhook Webhook `json:"updateWebhook"`
}

type DeleteWebhookResponse struct {
	Success bool `json:"deleteWebhook"`
}

func (p *clientImpl) GetWebhooks(ctx context.Context, params *WebhookParams) ([]Webhook, error) {
	query := `query webhooks($params: WebhookParams) { webhooks(params: $params) { id orgId org { name } url token description } }`

	variables := map[string]interface{}{
		"params": params,
	}

	var res GetWebhooksResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.Webhooks, nil
}

func (p *clientImpl) GetWebhook(ctx context.Context, id int) (*Webhook, error) {
	query := `query webhook($id: ID) { webhook(id: $id) { id url description token orgId org { name } }`

	variables := map[string]interface{}{
		"id": id,
	}

	var res GetWebhookResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Webhook, nil
}

func (p *clientImpl) CreateWebhook(ctx context.Context, webhook Webhook) (*Webhook, error) {
	query := `mutation createWebhook($input: WebhookInput!) {
  createWebhook(input: $input) {
    id
    url
    description
    token
    orgId
    org {
      id
      name
    }
    createdBy
    updatedBy
    createdAt
    updatedAt
  }
}
`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"orgId":         webhook.OrgID,
			"url":           webhook.URL,
			"token":         webhook.Token,
			"description":   webhook.Description,
			"generateToken": webhook.GenerateToken,
			"eventTypes":    webhook.EventTypes,
		},
	}

	var res CreateWebhookResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Webhook, nil
}

func (p *clientImpl) UpdateWebhook(ctx context.Context, webhook Webhook) (*Webhook, error) {

	if webhook.ID == 0 {
		return nil, errors.New("webhook id is required")
	}

	query := `mutation updateWebhook($input: WebhookInput!) { updateWebhook(input: $input) { id orgId org { name } url token description } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id": webhook.ID,
		},
	}

	if webhook.URL != "" {
		variables["input"].(map[string]interface{})["url"] = webhook.URL
	}

	if webhook.GenerateToken != nil {
		variables["input"].(map[string]interface{})["generateToken"] = webhook.GenerateToken
	}

	if webhook.Token != "" {
		variables["input"].(map[string]interface{})["token"] = webhook.Token
	}

	if webhook.Description != "" {
		variables["input"].(map[string]interface{})["description"] = webhook.Description
	}

	var res UpdateWebhookResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Webhook, nil
}

func (p *clientImpl) DeleteWebhook(ctx context.Context, id int) error {
	query := `mutation deleteWebhook($id: ID!) { deleteWebhook(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	var res DeleteWebhookResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return err
	}

	if !res.Success {
		return errors.New("failed to delete webhook")
	}

	return nil
}
