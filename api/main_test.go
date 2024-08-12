package api

import (
	"os"
	"testing"
	"time"

	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/arkarsg/splitapp/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.ServerConfig{
		Token: utils.TokenEnvs{
			SymmetricKey:   utils.RandomString(32),
			AccessDuration: time.Minute,
		},
	}
	server, err := NewServer(config, store)
	assert.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
