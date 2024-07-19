package utils

import (
	"testing"

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
	assert.Equal(t, "postgres://root:password@localhost:5432/split_db?sslmode=disable", src)
}
