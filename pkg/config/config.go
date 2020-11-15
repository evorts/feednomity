package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type App struct {
	Port              int           `yaml:"port"`
	HashSalt          string        `yaml:"hash_salt"`
	AESSalt           string        `yaml:"aes_salt"`
	SessionExpiration time.Duration `yaml:"session_expire"`
	CookieDomain      string        `yaml:"cookie_domain"`
	CookieSecure      int           `yaml:"cookie_secure"`
	Cors              struct {
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"cors"`
	TemplateDirectory string `yaml:"template_dir"`
	AssetDirectory    string `yaml:"asset_dir"`
	Logo              struct {
		FavIcon string `yaml:"favicon"`
		Url     string `yaml:"url"`
		Alt     string `yaml:"alt"`
	} `yaml:"logo"`
	Contact struct {
		Email   string   `yaml:"email"`
		Phone   []string `yaml:"phone"`
		Address string   `yaml:"address"`
	} `yaml:"contact"`
}

type Config struct {
	App App `yaml:"app"`
	DB  struct {
		Dsn                   string `yaml:"dsn"`
		MaxConnectionLifetime int64  `yaml:"max_connection_lifetime"`
		MaxIdleConnection     int64  `yaml:"max_idle_connection"`
		MaxOpenConnection     int64  `yaml:"max_open_connection"`
	} `yaml:"db"`
}

type config struct {
	filename []string
	data     *Config
}

type IManager interface {
	GetConfig() *Config
	Initiate() (IManager, error)
	Reload() error
}

func NewConfig(filename ...string) IManager {
	return &config{
		filename: filename,
		data:     nil,
	}
}

func (c *config) GetConfig() *Config {
	if c.data == nil {
		c.data, _ = c.read()
	}
	return c.data
}

func (c *config) GetApp() App {
	return c.data.App
}

func (c *config) Initiate() (IManager, error) {
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	c.data = data
	return c, nil
}

func (c *config) Reload() (err error) {
	c.data, err = c.read()
	return
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func (c *config) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (c *config) read() (*Config, error) {
	fName := ""
	for _, f := range c.filename {
		if c.fileExists(f) {
			fName = f
			break
		}
	}
	if len(fName) < 1 {
		return nil, errors.New("no configuration file found")
	}
	cfg, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(cfg, &config)
	if err != nil {
		return nil, err

	}
	return &config, nil
}
