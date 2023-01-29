package lib

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshal(t *testing.T) {
	input := HostsFile{
		Name:        "Name",
		Description: "Description",
		Content:     "Content",
		Url:         "Url",
		Enabled:     true,
	}
	expected := `---
Name: "Name"
Description: "Description"
Url: "Url"
---
Content`

	actual := Marshal(input)
	assert.Equal(t, expected, actual)
}

func TestUnmarshal(t *testing.T) {
	input := `---
Name: "Name"
Description: "Description"
Url: "Url"
---
Content`
	actual, err := Unmarshal(input)
	require.Nil(t, err)
	assert.Equal(t, "Name", actual.Name)
	assert.Equal(t, "Description", actual.Description)
	assert.Equal(t, "Url", actual.Url)
	assert.Equal(t, "Content", actual.Content)
}

func TestReplaceHostsFile(t *testing.T) {
	tests := []struct {
		content   string
		hostsFile []HostsFile
		expected  string
	}{
		{
			content: "",
			hostsFile: []HostsFile{
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
			content: `# before
########## START hosts-selector ##########
########## END   hosts-selector ##########
# after
`,
			hostsFile: []HostsFile{
				{
					Name:        "Name",
					Description: "Description",
					Content:     "Content",
					Url:         "",
					Enabled:     true,
				},
			},
			expected: `# before
########## START hosts-selector ##########
########################################
# Name
########################################
Content

########## END   hosts-selector ##########
# after
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
