package settings

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfigFileSuccess(t *testing.T) {
	configFile, err := BuildConfigFile("../../fixtures/settings/config.yml")

	require.Nil(t, err)

	assert.Equal(t,
		&ConfigFile{
			Version: "1",
			Org:     "foo",
			Services: []ConfigFileService{
				{
					Id:       "1",
					Name:     "foo bar service",
					Platform: "python",
					AuthKeys: []ConfigFileServiceAuthKey{
						{
							Key:       "key",
							Disabled:  false,
							ExpiredAt: 0,
						},
						{
							Key:       "expired_key",
							Disabled:  false,
							ExpiredAt: 946684800,
						},
						{
							Key:       "disabled_key",
							Disabled:  true,
							ExpiredAt: 0,
						},
					},
					Settings: ConfigFileServiceSettings{RateLimit: 1},
				},
			},
		}, configFile)
}

func TestParseConfigFileError(t *testing.T) {
	_, err := BuildConfigFile("../../fixtures/settings/config_invalid.yml")

	assert.True(t, errors.Is(err, ErrParseConfigFile))
}

func TestParseConfigFileErrorNotFound(t *testing.T) {
	_, err := BuildConfigFile("../../fixtures/settings/config_not_found.yml")

	assert.True(t, errors.Is(err, ErrConfigFileNotFound))
}
