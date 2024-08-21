package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	tdb := GetDevDbEnvs()
	s := GetServerEnvs()
	assert.NotEmpty(t, tdb)
	assert.NotEmpty(t, s)
}

func TestDbSource(t *testing.T) {
	src := GetDevDbSource()
	assert.Equal(t, "postgres://root:password@split-db:5432/split_db?sslmode=disable", src)
}

func TestTokenEnvs(t *testing.T) {
	envs := GetTokenEnvs()
	assert.Equal(t, time.Minute*15, envs.AccessDuration)
	assert.Len(t, envs.SymmetricKey, 32)
}
