package test

import (
	"testing"

	"github.com/sri-shubham/snipr/internal/config"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	confFile := "./config_test.yml"
	appConf, err := config.ParseConfig(confFile)
	require.Nil(t, err)
	require.NotNil(t, appConf.Redis)
	require.NotNil(t, appConf.Postgres)
	require.Equal(t, "localhost", appConf.Redis.Host)
	require.Equal(t, 6379, appConf.Redis.Port)
	require.Equal(t, 0, appConf.Redis.DB)
	require.Equal(t, 5, appConf.Redis.Timeout)

	require.Equal(t, "localhost", appConf.Postgres.Host)
	require.Equal(t, 5432, appConf.Postgres.Port)
	require.Equal(t, "snipr", appConf.Postgres.DB)
	require.Equal(t, "test", appConf.Postgres.User)
	require.Equal(t, "test", appConf.Postgres.Password)

	require.Equal(t, "snipr", appConf.Name)
	require.Equal(t, 8080, appConf.Port)
	require.Equal(t, "localhost", appConf.Host)

	require.NotNil(t, appConf.Shortener)
	require.Equal(t, 4, appConf.Shortener.MinLength)
	require.Equal(t, 7, appConf.Shortener.CustomMinLength)
	require.Equal(t, 16, appConf.Shortener.CustomMaxLength)

}
