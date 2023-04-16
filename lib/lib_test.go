package lib

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReplaceHostsFile(t *testing.T) {
	tests := []struct {
		content   string
		hostsFile []Hosts
		expected  string
	}{
		{
			content: "",
			hostsFile: []Hosts{
				{
					Name:        "Name",
					Description: "Description",
					Content:     "Content",
					Url:         "",
					Enabled:     true,
				},
			},
			expected: `
########## START hosts-selector ##########
########################################
# Name
########################################
Content

########## END   hosts-selector ##########
`,
		},
		{
			content: `before
########## START hosts-selector ##########
########## END   hosts-selector ##########
after
`,
			hostsFile: []Hosts{
				{
					Name:        "Name",
					Description: "Description",
					Content:     "Content",
					Url:         "",
					Enabled:     true,
				},
			},
			expected: `before
########## START hosts-selector ##########
########################################
# Name
########################################
Content

########## END   hosts-selector ##########
after
`,
		},
	}

	for _, tt := range tests {
		w := bytes.Buffer{}
		err := ReplaceHostsFile(tt.content, &w, tt.hostsFile)
		require.Nil(t, err)
		assert.Equal(t, tt.expected, w.String())
	}
}

func TestReplaceHostsFile2(t *testing.T) {
	actual := &bytes.Buffer{}
	hosts := HostsSet{
		Hosts{
			Name:        "hosts1",
			Description: "description1",
			Content:     "content1",
			Url:         "url1",
			Enabled:     true,
		},
		Hosts{
			Name:        "hosts2",
			Description: "description2",
			Content:     "content2",
			Url:         "url2",
			Enabled:     false,
		},
		Hosts{
			Name:        "hosts3",
			Description: "description3",
			Content:     "content3",
			Url:         "url3",
			Enabled:     true,
		},
	}

	input := "original hosts"
	expected := `original hosts
########## START hosts-selector ##########
########################################
# hosts1
########################################
content1
########################################
# hosts3
########################################
content3

########## END   hosts-selector ##########
`
	err := ReplaceHostsFile(input, actual, hosts)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.String())

	input = actual.String()
	actual = &bytes.Buffer{}
	err = ReplaceHostsFile("original hosts", actual, hosts)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.String())

}
