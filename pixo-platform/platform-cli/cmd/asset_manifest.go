package cmd

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"gopkg.in/yaml.v3"
	"os"
)

const ManifestFilename = "manifest.yaml"

type Asset struct {
	Name     string    `yaml:"name,omitempty"`
	Type     string    `yaml:"type,omitempty"`
	Versions []Version `yaml:"versions,omitempty"`
}

type Version struct {
	ID           int    `yaml:"id,omitempty"`
	Status       string `yaml:"status,omitempty"`
	LanguageCode string `yaml:"languageCode,omitempty"`
}

type Manifest struct {
	ModuleID           int     `yaml:"moduleId"`
	ModuleAbbreviation string  `yaml:"moduleAbbreviation,omitempty"`
	Assets             []Asset `yaml:"assets"`
}

func NewManifest() (*Manifest, error) {
	var manifest Manifest
	contents, err := os.ReadFile(ManifestFilename)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(contents, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (m Manifest) Save() error {
	manifestBytes, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	if err = os.WriteFile(ManifestFilename, manifestBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (m *Manifest) AddAsset(asset platform.Asset) {
	m.Assets = append(m.Assets, Asset{
		Name: asset.Name,
		Type: asset.Type,
	})
	for _, version := range asset.Versions {
		m.AddVersion(asset.Name, version)
	}
}

func (m *Manifest) GetAsset(name string) *Asset {
	for _, asset := range m.Assets {
		if asset.Name == name {
			return &asset
		}
	}
	return nil
}

func (m *Manifest) AddVersion(assetName string, version platform.AssetVersion) {
	asset := m.GetAsset(assetName)
	if asset == nil {
		return
	}

	asset.Versions = append(asset.Versions, Version{
		ID:           version.ID,
		Status:       version.Status,
		LanguageCode: version.LanguageCode,
	})

	m.SetAsset(assetName, asset)
}

func (m *Manifest) SetVersionToStatus(assetName string, versionID int, languageCode, startStatus, endStatus string) {
	asset := m.GetAsset(assetName)
	if asset == nil {
		return
	}
	for i, v := range asset.Versions {
		if v.ID == versionID {
			if v.Status == startStatus && (languageCode == "" || v.LanguageCode == languageCode) {
				asset.Versions[i].Status = endStatus
				m.SetAsset(assetName, asset)
			}
		}
	}
}

func (m *Manifest) SetAsset(name string, asset *Asset) {
	for i, a := range m.Assets {
		if a.Name == name {
			m.Assets[i] = *asset
			return
		}
	}
	m.Assets = append(m.Assets, *asset)
}
