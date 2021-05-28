package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type MailProvider map[string] string

func (m MailProvider) Get(key string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return ""
}

type MapMailProvider map[string]MailProvider

func (m MapMailProvider) Get(key string) MailProvider {
	if v, ok := m[key]; ok {
		return v
	}
	return MailProvider{}
}

type MemoryProvider struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type MapMemoryProvider map[string]MemoryProvider

func (m MapMemoryProvider) Get(key string) MemoryProvider {
	if v, ok := m[key]; ok {
		return v
	}
	return MemoryProvider{}
}

type App struct {
	Port               int    `yaml:"port"`
	PortApi            int    `yaml:"port_api"`
	BaseUrlWeb         string `yaml:"base_url_web"`
	BaseUrlApi         string `yaml:"base_url_api"`
	ReviewMailTemplate string `yaml:"review_mail_template"`
	HashSalt           string `yaml:"hash_salt"`
	AESSalt            string `yaml:"aes_salt"`
	HashExpire         int    `yaml:"hash_expire"`
	SessionExpiration  int64  `yaml:"session_expire"`
	CookieDomain       string `yaml:"cookie_domain"`
	CookieSecure       int    `yaml:"cookie_secure"`
	Cors               struct {
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"cors"`
	TemplateDirectory     string `yaml:"template_dir"`
	MailTemplateDirectory string `yaml:"mail_template_dir"`
	AssetDirectory        string `yaml:"asset_dir"`
	Logo                  struct {
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
	Jwe struct {
		Key    string `yaml:"key"`
		Expire int64  `yaml:"expire"`
	} `yaml:"jwe"`
	DB struct {
		Dsn                   string `yaml:"dsn"`
		MaxConnectionLifetime int64  `yaml:"max_connection_lifetime"`
		MaxIdleConnection     int64  `yaml:"max_idle_connection"`
		MaxOpenConnection     int64  `yaml:"max_open_connection"`
	} `yaml:"db"`
	Memory MapMemoryProvider `yaml:"memory"`
	Mailer struct {
		DailyLimit  int             `yaml:"daily_limit"`
		SenderName  string          `yaml:"sender_name"`
		SenderEmail string          `yaml:"sender_email"`
		Providers   MapMailProvider `yaml:"providers"`
	} `yaml:"mailer"`
	CronJobs struct {
		Blaster struct {
			Schedule  string `yaml:"schedule"`
			BatchRows int    `yaml:"batch_rows"`
		} `yaml:"blaster"`
	} `yaml:"cron_jobs"`
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

func (c *config) getPath(filename string) []string {
	p := make([]string, 0)
	wd, err := os.Getwd()
	if err == nil {
		p = append(p, path.Join(wd, filename))
	}
	usr, err := user.Current()
	if err == nil {
		p = append(p, path.Join(usr.HomeDir, filename))
	}
	return p
}

func (c *config) read() (*Config, error) {
	fName := ""
	for _, f := range c.filename {
		if c.fileExists(f) {
			fName = f
		}
		if len(fName) > 0 {
			break
		}
		paths := c.getPath(f)
		if len(paths) < 1 {
			continue
		}
		for _, ff := range paths {
			if !c.fileExists(ff) {
				continue
			}
			fName = ff
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
