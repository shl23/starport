package conf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/imdario/mergo"
)

var (
	// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
	ErrCouldntLocateConfig = errors.New("could not locate a config.yml in your chain. please follow the link for how-to: https://github.com/tendermint/starport/blob/develop/docs/1%20Introduction/4%20Configuration.md")

	// FileNames holds a list of appropriate names for the config file.
	FileNames = []string{"config.yml", "config.yaml"}

	// DefaultConf holds default configuration.
	DefaultConf = Config{
		Servers: Servers{
			RPCAddr:      "0.0.0.0:26657",
			P2PAddr:      "0.0.0.0:26656",
			ProfAddr:     "localhost:6060",
			GRPCAddr:     "0.0.0.0:9090",
			APIAddr:      "0.0.0.0:1317",
			FrontendAddr: "localhost:8080",
			DevUIAddr:    "localhost:12345",
		},
	}
)

// Config is the user given configuration to do additional setup
// during serve.
type Config struct {
	Accounts  []Account              `yaml:"accounts"`
	Validator Validator              `yaml:"validator"`
	Init      Init                   `yaml:"init"`
	Genesis   map[string]interface{} `yaml:"genesis"`
	Servers   Servers                `yaml:"servers"`
}

// AccountByName finds account by name.
func (c Config) AccountByName(name string) (acc Account, found bool) {
	for _, acc := range c.Accounts {
		if acc.Name == name {
			return acc, true
		}
	}
	return Account{}, false
}

// Account holds the options related to setting up Cosmos wallets.
type Account struct {
	Name     string   `yaml:"name"`
	Coins    []string `yaml:"coins,omitempty"`
	Mnemonic string   `yaml:"mnemonic,omitempty"`

	// The RPCAddress off the chain that account is issued at.
	RPCAddress string `yaml:"rpc_address,omitempty"`
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
}

// Init overwrites sdk configurations with given values.
type Init struct {
	// App overwrites appd's config/app.toml configs.
	App map[string]interface{} `yaml:"app"`

	// Config overwrites appd's config/config.toml configs.
	Config map[string]interface{} `yaml:"config"`

	// Home overwrites default home directory used for the app
	Home string `yaml:"home"`

	// CLIHome overwrites default CLI home directory used for launchpad app
	CLIHome string `yaml:"clihome"`
}

// Servers keeps configuration related to started servers.
type Servers struct {
	RPCAddr      string `yaml:"rpc-address"`
	P2PAddr      string `yaml:"p2p-address"`
	ProfAddr     string `yaml:"prof-address"`
	GRPCAddr     string `yaml:"grpc-address"`
	APIAddr      string `yaml:"api-address"`
	FrontendAddr string `yaml:"frontend-address"`
	DevUIAddr    string `yaml:"dev-ui-address"`
}

// Parse parses config.yml into UserConfig.
func Parse(r io.Reader) (Config, error) {
	var conf Config
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return conf, err
	}
	if err := mergo.Merge(&conf, DefaultConf); err != nil {
		return Config{}, err
	}
	return conf, validate(conf)
}

// ParseFile parses config.yml from the path.
func ParseFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, nil
	}
	defer file.Close()
	return Parse(file)
}

// validate validates user config.
func validate(conf Config) error {
	if len(conf.Accounts) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}
	if conf.Validator.Name == "" {
		return &ValidationError{"validator is required"}
	}
	return nil
}

// ValidationError is returned when a configuration is invalid.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config is not valid: %s", e.Message)
}

// Locate locates the path for the config file otherwise returns ErrCouldntLocateConfig.
func Locate(root string) (path string, err error) {
	for _, name := range FileNames {
		path = filepath.Join(root, name)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path, nil
		}
	}
	return "", ErrCouldntLocateConfig
}
