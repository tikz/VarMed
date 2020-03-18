package main

// Configuration holds constant parameters for the app
type Configuration struct {
	HTTPClient struct {
		UserAgent string `yaml:"user-agent"`
		Timeout   int    `yaml:"timeout"`
	} `yaml:"http-client"`

	HTTPServer struct {
		Port string `yaml:"port"`
	} `yaml:"http-server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}
