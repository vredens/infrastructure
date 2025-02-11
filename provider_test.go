package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTemplate(t *testing.T) {
	os.Setenv("INFRA_TEST_VAR", "test_value")
	provider, err := NewProvider(ProviderSettings{
		EnvName:       "test",
		SystemName:    "system",
		ComponentName: "comp",
	})
	assert.Nil(t, err)
	assert.Equal(t, "test", provider.RenderSecret("test"))
	assert.Equal(t, "{{ test }}", provider.RenderSecret("{{ test }}"))
	assert.Equal(t, "this is a test", provider.RenderSecret("this is a {{ .Environment }}"))

	// Test invalid template
	result, err := provider.RenderSecrets("{{ .InvalidProp }}")
	assert.NotNil(t, err)
	assert.Empty(t, result)

	// Test environment variable rendering
	assert.Equal(t, "prefix_test_value_suffix", provider.RenderSecret("prefix_{{ .Env.INFRA_TEST_VAR }}_suffix"))

	// Test empty template
	assert.Equal(t, "", provider.RenderSecret(""))

	// Test component name rendering
	assert.Equal(t, "comp", provider.RenderSecret("{{ .Component }}"))
}

func TestProviderEnvironmentInfo(t *testing.T) {
	provider, err := NewProvider(ProviderSettings{
		EnvName:       "test",
		SystemName:    "system",
		ComponentName: "comp",
	})
	assert.Nil(t, err)

	assert.Equal(t, "test", provider.Environment())
	assert.Equal(t, "system", provider.SystemName())
	assert.Equal(t, "comp", provider.ComponentName())
}

func TestProviderWithCustomSettings(t *testing.T) {
	settings := ProviderSettings{
		EnvName:       "test",
		SystemName:    "customsys",
		ComponentName: "customcomp",
		CertFolders: []string{
			"testdata/custom/certs",
		},
		AppConfigFolders: []string{
			"testdata/custom/config",
		},
		InfraConfigFolders: []string{
			"testdata/custom/infra",
		},
	}

	provider, err := NewProvider(settings)
	assert.Nil(t, err)
	assert.Equal(t, "test", provider.Environment())
	assert.Equal(t, "customsys", provider.SystemName())
	assert.Equal(t, "customcomp", provider.ComponentName())
}

func TestProviderInvalidSettings(t *testing.T) {
	testCases := []struct {
		name     string
		settings ProviderSettings
		errMsg   string
	}{
		{
			name: "empty all",
			settings: ProviderSettings{
				EnvName:       "",
				SystemName:    "",
				ComponentName: "",
			},
			errMsg: "could not determine environment name",
		},
		{
			name: "empty system and component",
			settings: ProviderSettings{
				EnvName:       "test",
				SystemName:    "",
				ComponentName: "",
			},
			errMsg: "could not determine system name from environment",
		},
		{
			name: "empty component only",
			settings: ProviderSettings{
				EnvName:       "test",
				SystemName:    "sys",
				ComponentName: "",
			},
			errMsg: "could not determine component name from environment",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.settings.Validate()
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func BenchmarkRenderTemplate(b *testing.B) {
	os.Setenv("INFRA_TEST_VAR", "test_value")
	provider, err := NewProvider(ProviderSettings{
		EnvName:       "test",
		SystemName:    "system",
		ComponentName: "comp",
	})
	assert.Nil(b, err)

	benchCases := []struct {
		name     string
		template string
	}{
		{
			name:     "simple text",
			template: "simple text without template",
		},
		{
			name:     "with env var",
			template: "this is a {{ .Env.INFRA_TEST_VAR }}",
		},
		{
			name:     "with component",
			template: "component: {{ .Component }}",
		},
		{
			name:     "complex template",
			template: "env: {{ .Environment }}, system: {{ .System }}, component: {{ .Component }}",
		},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := provider.RenderSecret(bc.template)
				if result == "" {
					b.FailNow()
				}
			}
		})
	}
}

func TestProviderLoadConfig(t *testing.T) {
	settings := ProviderSettings{
		EnvName:       "test",
		SystemName:    "sys",
		ComponentName: "cmp",
	}
	t.Run("incomplete struct", func(t *testing.T) {
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			Repo string `json:"repo-1"`
		}
		var cfg myTestConfig
		assert.NotNil(t, provider.LoadConfig("app", &cfg))
		provider.Locator()
	})

	t.Run("complete struct", func(t *testing.T) {
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct {
				Resource string `json:"arn"`
				Params   struct {
					Timeout            int `json:"timeout"`
					MaxIdleConnections int `json:"max_idle_connections"`
				} `json:"params"`
			} `json:"repo-1"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfig("app", &cfg))
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo01.Resource)
		assert.Equal(t, 10, cfg.SampleRepo01.Params.Timeout)
		provider.Locator()
	})

	t.Run("load specific config for env", func(t *testing.T) {
		provider, err := NewProvider(ProviderSettings{
			EnvName:       "testenv",
			SystemName:    "sys",
			ComponentName: "cmp",
		})
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct { // partial override in the env config
				Resource string `json:"arn"`
				Params   struct {
					Timeout            int `json:"timeout"`
					MaxIdleConnections int `json:"max_idle_connections"`
				} `json:"params"`
			} `json:"repo-1"`
			SampleRepo02 struct { // full override in the env config
				Resource string `json:"arn"`
				Params   struct {
					Timeout            int `json:"timeout"`
					MaxIdleConnections int `json:"max_idle_connections"`
				} `json:"params"`
			} `json:"repo-2"`
			SampleRepo03 struct { // not in the env config
				Resource string `json:"arn"`
			} `json:"repo-6"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfig("app", &cfg))
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo01.Resource)
		assert.Equal(t, 15, cfg.SampleRepo01.Params.Timeout)
		assert.Equal(t, "arn://storage/elasticsearch/sample-2", cfg.SampleRepo02.Resource)
		assert.Equal(t, 20, cfg.SampleRepo02.Params.Timeout)
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo03.Resource)
		provider.Locator()
	})

	t.Run("load specific config for env without global", func(t *testing.T) {
		provider, err := NewProvider(ProviderSettings{
			EnvName:       "testenv",
			SystemName:    "sys",
			ComponentName: "cmp",
		})
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo struct { // partial override in the env config
				Resource string `json:"arn"`
				Params   struct {
					Timeout            int `json:"timeout"`
					MaxIdleConnections int `json:"max_idle_connections"`
				} `json:"params"`
			} `json:"repo"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfig("env-app", &cfg))
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo.Resource)
		assert.Equal(t, 15, cfg.SampleRepo.Params.Timeout)
		provider.Locator()
	})

	t.Run("render config with secrets", func(t *testing.T) {
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct {
				Resource string `json:"arn"`
				Params   struct {
					Password  string `json:"password"`
					Password2 string `json:"password2"`
					Password3 string `json:"password3"`
				} `json:"params"`
			} `json:"repo-6"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfig("app", &cfg))
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo01.Resource)
		assert.Equal(t, "never-use-this-pass-test-sys-cmp", cfg.SampleRepo01.Params.Password)
		assert.Equal(t, "", cfg.SampleRepo01.Params.Password3)
		provider.Locator()
	})

	t.Run("failed to render config with secrets", func(t *testing.T) {
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct {
				Resource string `json:"arn"`
				Params   struct {
					Password  string `json:"password"`
					Password2 string `json:"password2"`
					Password3 string `json:"password3"`
				} `json:"params"`
			} `json:"repo-6"`
		}
		var cfg myTestConfig
		assert.NotNil(t, provider.LoadConfig("app-invalid", &cfg))
		provider.Locator()
	})

	t.Run("render config with Infra secrets", func(t *testing.T) {
		os.Setenv("INFRA_KAFKA_USERNAME", "useruser")
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct {
				Resource string `json:"arn"`
				Params   struct {
					UserName  string `json:"username"`
					UserName2 string `json:"username2"`
					UserName3 string `json:"username3"`
				} `json:"params"`
			} `json:"repo-7"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfig("app", &cfg))
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo01.Resource)
		assert.Equal(t, "useruser", cfg.SampleRepo01.Params.UserName)
		assert.Empty(t, cfg.SampleRepo01.Params.UserName2)
		assert.Equal(t, "this is a useruser", cfg.SampleRepo01.Params.UserName3)
		provider.Locator()
	})

	t.Run("render config with Infra secrets with special character", func(t *testing.T) {
		os.Setenv("INFRA_KAFKA_USERNAME", "user=user")
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct {
				Resource string `json:"arn"`
				Params   struct {
					UserName  string `json:"username"`
					UserName2 string `json:"username2"`
					UserName3 string `json:"username3"`
				} `json:"params"`
			} `json:"repo-7"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfig("app", &cfg))
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo01.Resource)
		assert.Equal(t, "useruser", cfg.SampleRepo01.Params.UserName)
		assert.Empty(t, cfg.SampleRepo01.Params.UserName2)
		assert.Equal(t, "this is a user=user", cfg.SampleRepo01.Params.UserName3)
		provider.Locator()
	})

	t.Run("render config with Infra secrets with special character", func(t *testing.T) {
		os.Setenv("INFRA_KAFKA_USERNAME", "user=user")
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct {
				Resource string `json:"arn"`
				Params   struct {
					UserName  string `json:"username"`
					UserName2 string `json:"username2"`
					UserName3 string `json:"username3"`
				} `json:"params"`
			} `json:"repo-7"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfig("app", &cfg))
		assert.Equal(t, "arn://storage/elasticsearch/sample-1", cfg.SampleRepo01.Resource)
		assert.Equal(t, "useruser", cfg.SampleRepo01.Params.UserName)
		assert.Empty(t, cfg.SampleRepo01.Params.UserName2)
		assert.Equal(t, "this is a user=user", cfg.SampleRepo01.Params.UserName3)
		provider.Locator()
	})

	t.Run("load_config_from_file_with_Infra_secrets", func(t *testing.T) {
		os.Setenv("INFRA_KAFKA_USERNAME", "user")
		os.Setenv("INFRA_KAFKA_PASSWORD", "pAss=Word")
		provider, err := NewProvider(settings)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		type myTestConfig struct {
			SampleRepo01 struct {
				Resource string `json:"arn"`
			} `json:"repo"`
		}
		var cfg myTestConfig
		assert.Nil(t, provider.LoadConfigFromFile("./testdata/config/from-file.json", &cfg))
		assert.Equal(t, "arn://messaging/kafka/clusters/c1", cfg.SampleRepo01.Resource)
		assert.Len(t, provider.infraConfig.Messaging.Kafka.Clusters, 1)
		cluster := provider.infraConfig.Messaging.Kafka.Clusters["c1"]
		assert.Equal(t, cluster.Username, "user")
		assert.Equal(t, cluster.Password, "pAss=Word")
	})
}
