package lib

import (
	"errors"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
	"strings"
)

type HostsFile struct {
	Name        string
	Description string
	Content     string
	Url         string
	Enabled     bool
}

type Config struct {
	HostsFilePath string      `toml:"hostsFilePath"`
	Hosts         []HostsFile `toml:"hosts"`
}

func (c *Config) SearchHostsFileName(name string) int {
	for i := range c.Hosts {
		if name == c.Hosts[i].Name {
			return i
		}
	}
	return -1
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

func ReplaceHostsFile(content string, w io.Writer, hosts []HostsFile) error {

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

func Marshal(data HostsFile) string {
	frontmatter := "---\n" +
		"Name: \"" + data.Name + "\"\n" +
		"Description: \"" + data.Description + "\"\n" +
		"Url: \"" + data.Url + "\"\n" +
		"---\n"

	return frontmatter + data.Content
}

func Unmarshal(data string) (HostsFile, error) {
	if !strings.HasPrefix(data, FrontMatterStart) {
		return HostsFile{}, errors.New("format error")
	}
	be := strings.Index(data, FrontMatterEnd)
	lines := strings.Split(data[len(FrontMatterStart):be], "\n")
	meta := make(map[string]string)
	for _, l := range lines {
		elements := strings.SplitN(l, FrontMatterDelimiter, 2)
		if len(elements) != 2 {
			return HostsFile{}, errors.New("format error")
		}
		key := strings.Trim(elements[0], " \t\r\n")
		value := strings.Trim(elements[1], " \t\r\n")
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = value[1 : len(value)-1]
		}
		meta[key] = value
	}

	result := HostsFile{}
	if val, ok := meta["Name"]; ok {
		result.Name = val
	}
	if val, ok := meta["Description"]; ok {
		result.Description = val
	}
	if val, ok := meta["Url"]; ok {
		result.Url = val
	}
	result.Content = data[be+len(FrontMatterEnd):]
	return result, nil
}
