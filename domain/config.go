package domain

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name         string   `yaml:"name"`
	Provider     string   `yaml:"provider"`
	Token        string   `yaml:"token"`
	URL          string   `yaml:"url"`
	Author       string   `yaml:"author"`
	Repositories []string `yaml:"repositories"`
	Owner        string   `yaml:"owner"`
}
