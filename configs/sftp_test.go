package configs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vredens/infrastructure"
	"github.com/vredens/infrastructure/configs"
)

func TestSFTP(t *testing.T) {
	provider, err := infrastructure.NewProvider(infrastructure.ProviderSettings{
		EnvName:       "sftp-tests",
		SystemName:    "tests",
		ComponentName: "test",
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Run("user-pass", func(t *testing.T) {
		var cfg struct {
			Test configs.SFTP `json:"valid-1"`
		}
		if !assert.NoError(t, provider.LoadConfig("sftp", &cfg)) {
			t.FailNow()
		}
		assert.NoError(t, cfg.Test.Bootstrap(provider))
		assert.Equal(t, "localhost", cfg.Test.Resource().Host)
		assert.Equal(t, 22, cfg.Test.Resource().Port)
		assert.Equal(t, "user", cfg.Test.Resource().User)
		assert.Equal(t, "pass", cfg.Test.Resource().Pass)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Value)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Path)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Passphrase)
	})

	t.Run("private-key", func(t *testing.T) {
		var cfg struct {
			Test configs.SFTP `json:"valid-2"`
		}
		if !assert.NoError(t, provider.LoadConfig("sftp", &cfg)) {
			t.FailNow()
		}
		assert.NoError(t, cfg.Test.Bootstrap(provider))
		assert.Equal(t, "localhost", cfg.Test.Resource().Host)
		assert.Equal(t, 22, cfg.Test.Resource().Port)
		assert.Equal(t, "user", cfg.Test.Resource().User)
		assert.Equal(t, "", cfg.Test.Resource().Pass)
		assert.Equal(t, "private_key", cfg.Test.Resource().PrivateKey.Value)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Path)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Passphrase)
	})

	t.Run("private-key-file", func(t *testing.T) {
		var cfg struct {
			Test configs.SFTP `json:"valid-3"`
		}
		if !assert.NoError(t, provider.LoadConfig("sftp", &cfg)) {
			t.FailNow()
		}
		assert.NoError(t, cfg.Test.Bootstrap(provider))
		assert.Equal(t, "localhost", cfg.Test.Resource().Host)
		assert.Equal(t, 22, cfg.Test.Resource().Port)
		assert.Equal(t, "user", cfg.Test.Resource().User)
		assert.Equal(t, "", cfg.Test.Resource().Pass)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Value)
		assert.Equal(t, "/path/to/private_key", cfg.Test.Resource().PrivateKey.Path)
		assert.Equal(t, "passphrase", cfg.Test.Resource().PrivateKey.Passphrase)
	})

	t.Run("private-key-with-passphrase", func(t *testing.T) {
		var cfg struct {
			Test configs.SFTP `json:"valid-4"`
		}
		if !assert.NoError(t, provider.LoadConfig("sftp", &cfg)) {
			t.FailNow()
		}
		assert.NoError(t, cfg.Test.Bootstrap(provider))
		assert.Equal(t, "localhost", cfg.Test.Resource().Host)
		assert.Equal(t, 22, cfg.Test.Resource().Port)
		assert.Equal(t, "user", cfg.Test.Resource().User)
		assert.Equal(t, "", cfg.Test.Resource().Pass)
		assert.Equal(t, "private_key", cfg.Test.Resource().PrivateKey.Value)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Path)
		assert.Equal(t, "passphrase", cfg.Test.Resource().PrivateKey.Passphrase)
	})

	t.Run("dual-auth", func(t *testing.T) {
		var cfg struct {
			Test configs.SFTP `json:"valid-5"`
		}
		if !assert.NoError(t, provider.LoadConfig("sftp", &cfg)) {
			t.FailNow()
		}
		assert.NoError(t, cfg.Test.Bootstrap(provider))
		assert.Equal(t, "localhost", cfg.Test.Resource().Host)
		assert.Equal(t, 22, cfg.Test.Resource().Port)
		assert.Equal(t, "user", cfg.Test.Resource().User)
		assert.Equal(t, "pass", cfg.Test.Resource().Pass)
		assert.Equal(t, "private_key", cfg.Test.Resource().PrivateKey.Value)
		assert.Equal(t, "", cfg.Test.Resource().PrivateKey.Path)
		assert.Equal(t, "passphrase", cfg.Test.Resource().PrivateKey.Passphrase)
	})

	t.Run("invalid-1", func(t *testing.T) {
		var cfg struct {
			Test configs.SFTP `json:"invalid-1"`
		}
		if !assert.NoError(t, provider.LoadConfig("sftp", &cfg)) {
			t.FailNow()
		}
		assert.Error(t, cfg.Test.Bootstrap(provider))
	})

	t.Run("invalid-2", func(t *testing.T) {
		var cfg struct {
			Test configs.SFTP `json:"invalid-2"`
		}
		if !assert.NoError(t, provider.LoadConfig("sftp", &cfg)) {
			t.FailNow()
		}
		assert.Error(t, cfg.Test.Bootstrap(provider))
	})
}
