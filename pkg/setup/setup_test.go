package setup

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConfigFromBytes_ExampleConfigIsCorrect(t *testing.T) {
	configBytes, err := os.ReadFile("../../config.example.yaml")
	if err != nil {
		assert.NoError(t, err, "cannot read config file: "+err.Error())
		return
	}

	config, err := ParseConfigFromBytes(configBytes)

	assert.NoError(t, err)
	assert.NotEmpty(t, config.LogLevel)
	assert.NotEmpty(t, config.LogFormat)
	assert.NotEmpty(t, config.Route)
	assert.NotEmpty(t, config.Route.Routes)
	assert.Equal(t, 4, len(config.Route.Routes))
	assert.NotEmpty(t, config.Receivers)
	assert.Equal(t, 11, len(config.Receivers))
}

func Test_ParseConfigFromBytes_NoErrors(t *testing.T) {
	configBytes := []byte(`
logLevel: info
logFormat: json
`)

	config, err := ParseConfigFromBytes(configBytes)

	assert.NoError(t, err)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "json", config.LogFormat)
}

func Test_ParseConfigFromBytes_ErrorWhenCurlyBracesNotEscaped(t *testing.T) {
	configBytes := []byte(`
logLevel: {{info}}
logFormat: json
`)

	config, err := ParseConfigFromBytes(configBytes)

	expectedErrorLine := ">  2 | logLevel: {{info}}"
	expectedErrorSuggestion := "Need to wrap values with special characters in quotes"
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorLine)
	assert.Contains(t, err.Error(), expectedErrorSuggestion)
	assert.Equal(t, "", config.LogLevel)
	assert.Equal(t, "", config.LogFormat)
}

func Test_ParseConfigFromBytes_OkWhenCurlyBracesEscaped(t *testing.T) {
	configBytes := []byte(`
logLevel: "{{info}}"
logFormat: json
`)

	config, err := ParseConfigFromBytes(configBytes)

	assert.Nil(t, err)
	assert.Equal(t, "{{info}}", config.LogLevel)
	assert.Equal(t, "json", config.LogFormat)
}

func Test_ParseConfigFromBytes_ErrorErrorNotWithCurlyBraces(t *testing.T) {
	configBytes := []byte(`
logLevelNotYAMLErrorError
logFormat: json
`)

	config, err := ParseConfigFromBytes(configBytes)

	expectedErrorLine := ">  2 | logLevelNotYAMLErrorError"
	expectedErrorSuggestion := "Need to wrap values with special characters in quotes"
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorLine)
	assert.NotContains(t, err.Error(), expectedErrorSuggestion)
	assert.Equal(t, "", config.LogLevel)
	assert.Equal(t, "", config.LogFormat)
}
