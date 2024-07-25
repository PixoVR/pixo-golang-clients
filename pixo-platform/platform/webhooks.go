package platform

import (
	"context"
	"encoding/json"
	"errors"
)

type Webhook struct {
	ID            int    `json:"id,omitempty"`
	URL           string `json:"url,omitempty"`
	Description   string `json:"description,omitempty"`
	Token         string `json:"token,omitempty"`
	GenerateToken *bool  `json:"generateToken,omitempty"`
	OrgID         int    `json:"orgId,omitempty"`
	Org           *Org   `json:"org,omitempty"`
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

func (g *PlatformClient) GetWebhooks(ctx context.Context, params *WebhookParams) ([]Webhook, error) {
	query := `query webhooks($params: WebhookParams) { webhooks(params: $params) { id orgId org { name } url token description } }`

	variables := map[string]interface{}{
		"params": params,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var webhooksResponse GetWebhooksResponse
	if err = json.Unmarshal(res, &webhooksResponse); err != nil {
		return nil, err

	}

	return webhooksResponse.Webhooks, nil
}

func (g *PlatformClient) GetWebhook(ctx context.Context, id int) (*Webhook, error) {
	query := `query webhook($id: ID) { webhook(id: $id) { id url description token orgId org { name } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var webhookResponse GetWebhookResponse
	if err = json.Unmarshal(res, &webhookResponse); err != nil {
		return nil, err
	}

	return &webhookResponse.Webhook, nil
}

func (g *PlatformClient) CreateWebhook(ctx context.Context, webhook Webhook) (*Webhook, error) {
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
			"orgId":       webhook.OrgID,
			"url":         webhook.URL,
			"token":       webhook.Token,
			"description": webhook.Description,
		},
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var webhookResponse CreateWebhookResponse
	if err = json.Unmarshal(res, &webhookResponse); err != nil {
		return nil, err
	}

	return &webhookResponse.Webhook, nil
}

func (g *PlatformClient) UpdateWebhook(ctx context.Context, webhook Webhook) (*Webhook, error) {

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

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var updateWebhookResponse UpdateWebhookResponse
	if err = json.Unmarshal(res, &updateWebhookResponse); err != nil {
		return nil, err
	}

	return &updateWebhookResponse.Webhook, nil
}

func (g *PlatformClient) DeleteWebhook(ctx context.Context, id int) error {
	query := `mutation deleteWebhook($id: ID!) { deleteWebhook(id: $id) }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return err
	}

	var deleteResponse DeleteWebhookResponse
	if err = json.Unmarshal(res, &deleteResponse); err != nil {
		return err
	}

	if !deleteResponse.Success {
		return errors.New("failed to delete webhook")
	}

	return nil
}
