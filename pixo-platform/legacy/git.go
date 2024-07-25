package legacy

type GitConfig struct {
	Provider string `json:"provider,omitempty"`
	OrgName  string `json:"orgName,omitempty"`
	RepoName string `json:"repoName,omitempty"`
}
