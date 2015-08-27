package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrDefault(t *testing.T) {
	envKey, envVal := "FOO", "someValue"
	os.Setenv(envKey, envVal)
	assert.Equal(t, envVal, GetOrDefault(envKey, ""))

	envKey = "BAR"
	defVal := "defaultValue"
	assert.Equal(t, defVal, GetOrDefault(envKey, defVal))
}
