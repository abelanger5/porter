package models

// IndexYAML represents a chart repo's index.yaml
type IndexYAML struct {
	APIVersion string                    `yaml:"apiVersion"`
	Generated  string                    `yaml:"generated"`
	Entries    map[interface{}]ChartYAML `yaml:"entries"`
}

// ChartYAML represents the data for chart in index.yaml
type ChartYAML []struct {
	APIVersion  string   `yaml:"apiVersion"`
	AppVersion  string   `yaml:"appVersion"`
	Created     string   `yaml:"created"`
	Description string   `yaml:"description"`
	Digest      string   `yaml:"digest"`
	Icon        string   `yaml:"icon"`
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	Urls        []string `yaml:"urls"`
	Version     string   `yaml:"version"`
}

// PorterChart represents a bundled Porter template
type PorterChart struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Form        FormYAML `json:"form"`
	Markdown    string   `json:"markdown"`
}

// FormYAML represents a chart's values.yaml form abstraction
type FormYAML struct {
	Name        string   `yaml:"name" json:"name"`
	Icon        string   `yaml:"icon" json:"icon"`
	Description string   `yaml:"description" json:"description"`
	Tags        []string `yaml:"tags" json:"tags"`
	Tabs        []struct {
		Name     string `yaml:"name" json:"name"`
		Label    string `yaml:"label" json:"label"`
		Sections []struct {
			Name     string `yaml:"name" json:"name"`
			ShowIf   string `yaml:"show_if" json:"show_if"`
			Contents []struct {
				Type     string `yaml:"type" json:"type"`
				Label    string `yaml:"label" json:"label"`
				Name     string `yaml:"name,omitempty" json:"name,omitempty"`
				Variable string `yaml:"variable,omitempty" json:"variable,omitempty"`
				Settings struct {
					Default interface{} `yaml:"default,omitempty" json:"default,omitempty"`
					Unit    interface{} `yaml:"unit,omitempty" json:"unit,omitempty"`
				} `yaml:"settings,omitempty" json:"settings,omitempty"`
			} `yaml:"contents" json:"contents,omitempty"`
		} `yaml:"sections" json:"sections,omitempty"`
	} `yaml:"tabs" json:"tabs,omitempty"`
}
