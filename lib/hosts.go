package lib

import (
	"errors"
	"strings"
)

type Hosts struct {
	Name        string
	Description string
	Content     string
	Url         string
	Enabled     bool
}

func Marshal(data Hosts) string {
	frontmatter := "---\n" +
		"Name: \"" + data.Name + "\"\n" +
		"Description: \"" + data.Description + "\"\n" +
		"Url: \"" + data.Url + "\"\n" +
		"---\n"

	return frontmatter + data.Content
}

func Unmarshal(data string) (Hosts, error) {
	if !strings.HasPrefix(data, FrontMatterStart) {
		return Hosts{}, errors.New("format error")
	}
	be := strings.Index(data, FrontMatterEnd)
	lines := strings.Split(data[len(FrontMatterStart):be], "\n")
	meta := make(map[string]string)
	for _, l := range lines {
		elements := strings.SplitN(l, FrontMatterDelimiter, 2)
		if len(elements) != 2 {
			return Hosts{}, errors.New("format error")
		}
		key := strings.Trim(elements[0], " \t\r\n")
		value := strings.Trim(elements[1], " \t\r\n")
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = value[1 : len(value)-1]
		}
		meta[key] = value
	}

	result := Hosts{}
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

type HostsSet []Hosts

func (s *HostsSet) SearchByName(name string) *Hosts {
	for i := range *s {
		if name == (*s)[i].Name {
			return &(*s)[i]
		}
	}
	return nil
}
