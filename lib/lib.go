package lib

import (
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
	"strings"
)

type Config struct {
	HostsFilePath string   `toml:"hostsFilePath"`
	Hosts         HostsSet `toml:"hosts"`
}

const DefaultHostsFilePath = "/etc/hosts"

func LoadConfig(filepath string) (*Config, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := toml.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	if config.HostsFilePath == "" {
		config.HostsFilePath = DefaultHostsFilePath
	}

	return &config, err
}

func SaveConfig(filepath string, config *Config) error {
	b, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, b, 0666)
}

const (
	BoundaryStart = "\n########## START hosts-selector ##########\n"
	BoundaryEnd   = "\n########## END   hosts-selector ##########\n"
)

func ReplaceHostsFile(content string, w io.Writer, hosts HostsSet) error {

	sp := strings.Index(content, BoundaryStart)
	ep := strings.Index(content, BoundaryEnd)

	if sp != -1 && ep != -1 && sp < ep {
		io.WriteString(w, content[:sp+len(BoundaryStart)])
	} else {
		io.WriteString(w, content)
		io.WriteString(w, BoundaryStart)
	}

	for _, h := range hosts {
		if h.Enabled {
			io.WriteString(w, "########################################\n# ")
			io.WriteString(w, h.Name)
			io.WriteString(w, "\n########################################\n")
			io.WriteString(w, h.Content)
			io.WriteString(w, "\n")
		}
	}

	if sp != -1 && ep != -1 && sp < ep {
		io.WriteString(w, content[ep:])
	} else {
		io.WriteString(w, BoundaryEnd)
	}

	return nil
}

const (
	FrontMatterStart     = "---\n"
	FrontMatterEnd       = "\n---\n"
	FrontMatterDelimiter = ":"
)
