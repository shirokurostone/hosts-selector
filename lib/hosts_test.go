package lib

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshal(t *testing.T) {
	input := Hosts{
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
