package configs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vredens/infrastructure"
	"github.com/vredens/infrastructure/configs"
)

func TestPostgresResource(t *testing.T) {
	os.Setenv("PSQL_USER", "username")
	os.Setenv("PSQL_PASS", "password")
	os.Setenv("INFRA_PSQL_USER", "user")
	os.Setenv("INFRA_PSQL_PASS", "pass")
	provider, err := infrastructure.NewProvider(infrastructure.ProviderSettings{
		EnvName:       "pg-tests",
		SystemName:    "tests",
		ComponentName: "test",
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Run("valid-1", func(t *testing.T) {
		var cfg struct {
			Test configs.Postgres `json:"valid-1"`
		}
		if !assert.NoError(t, provider.LoadConfig("pg", &cfg)) {
			t.FailNow()
		}
		assert.NoError(t, cfg.Test.Bootstrap(provider))
		assert.Equal(t, "localhost", cfg.Test.Resource().Host)
		assert.Equal(t, uint16(5432), cfg.Test.Resource().Port)
		assert.Equal(t, "username", cfg.Test.Resource().User)
		assert.Equal(t, "password", cfg.Test.Resource().Password)
	})

	t.Run("valid-2", func(t *testing.T) {
		var cfg struct {
			Test configs.Postgres `json:"valid-2"`
		}
		if !assert.NoError(t, provider.LoadConfig("pg", &cfg)) {
			t.FailNow()
		}
		assert.NoError(t, cfg.Test.Bootstrap(provider))
		assert.Equal(t, "localhost", cfg.Test.Resource().Host)
		assert.Equal(t, uint16(5432), cfg.Test.Resource().Port)
		assert.Equal(t, "user", cfg.Test.Resource().User)
		assert.Equal(t, "pass", cfg.Test.Resource().Password)
	})
}
