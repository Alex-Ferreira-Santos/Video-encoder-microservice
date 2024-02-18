package utils_test

import (
	"testing"

	"github.com/alex-ferreira-santos/encoder/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestIsJSON(t *testing.T) {
	json := `{
		"id": "1a45f0fb-d0f5-4cd6-a212-6e626b08af8c",
		"file_path": "Fly me to the moon.mp4",
		"status": "pending"
	}`

	err := utils.IsJson(json)

	require.Nil(t, err)

	json = `alex`

	err = utils.IsJson(json)

	require.Error(t, err)
}