package platform

import (
	"context"
)

const (
	LifecycleDevelopment = "Development"
	LifecycleQA          = "QA"
	LifecycleReleased    = "Released"
	LifecycleArchived    = "Archived"
)

type VersionLifecycle struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type VersionLifecyclesResponse struct {
	VersionLifecycles []VersionLifecycle `json:"versionLifecycles"`
}

func (p *clientImpl) GetVersionLifecycles(ctx context.Context) ([]VersionLifecycle, error) {
	query := `query versionLifecycles { versionLifecycles { id name } }`

	var res VersionLifecyclesResponse
	if err := p.Exec(ctx, query, &res, nil); err != nil {
		return nil, err
	}

	return res.VersionLifecycles, nil
}
