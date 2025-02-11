package configs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vredens/infrastructure"
	"github.com/vredens/infrastructure/configs"
)

func TestKafkaResources(t *testing.T) {
	provider, err := infrastructure.NewProvider(infrastructure.ProviderSettings{
		EnvName:       "kafka-tests",
		SystemName:    "tests",
		ComponentName: "test",
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Run("consumer", func(t *testing.T) {
		var cfg struct {
			Kafka struct {
				Test configs.KafkaConsumer `json:"test-1"`
			}
		}
		if !assert.NoError(t, provider.LoadConfig("kafka", &cfg)) {
			t.FailNow()
		}
		if !assert.NoError(t, cfg.Kafka.Test.Bootstrap(provider)) {
			t.FailNow()
		}
		assert.Equal(t, "my.group", cfg.Kafka.Test.Group)
		assert.Equal(t, "gpa-my.group-gsa", cfg.Kafka.Test.GroupNameFor(cfg.Kafka.Test.Group))
		assert.Equal(t, "gpa-my.group-gsa", cfg.Kafka.Test.GroupName())
		assert.Equal(t, "my-topic", cfg.Kafka.Test.Topic)
		assert.Equal(t, "tpa-my-topic-tsa", cfg.Kafka.Test.TopicNameFor(cfg.Kafka.Test.Topic))
		assert.Equal(t, "tpa-my-topic-tsa", cfg.Kafka.Test.TopicName())
	})

	t.Run("consumer/topic/translation", func(t *testing.T) {
		var cfg struct {
			Kafka struct {
				Test configs.KafkaConsumer `json:"test-5"`
			}
		}
		if !assert.NoError(t, provider.LoadConfig("kafka", &cfg)) {
			t.FailNow()
		}
		if !assert.NoError(t, cfg.Kafka.Test.Bootstrap(provider)) {
			t.FailNow()
		}
		assert.Equal(t, "our-topic", cfg.Kafka.Test.TopicNameFor("my-topic"))
	})

	t.Run("consumer/deprecated", func(t *testing.T) {
		var cfg struct {
			Kafka struct {
				Test configs.KafkaConsumer `json:"test-2"`
			}
		}
		if !assert.NoError(t, provider.LoadConfig("kafka", &cfg)) {
			t.FailNow()
		}
		assert.Error(t, cfg.Kafka.Test.Bootstrap(provider))
	})

	t.Run("producer", func(t *testing.T) {
		var cfg struct {
			Kafka struct {
				Test configs.KafkaProducer `json:"test-3"`
			}
		}
		if !assert.NoError(t, provider.LoadConfig("kafka", &cfg)) {
			t.FailNow()
		}
		if !assert.NoError(t, cfg.Kafka.Test.Bootstrap(provider)) {
			t.FailNow()
		}
		assert.Equal(t, "my-topic", cfg.Kafka.Test.Topic)
		assert.Equal(t, "tpa-my-topic-tsa", cfg.Kafka.Test.TopicNameFor(cfg.Kafka.Test.Topic))
		assert.Equal(t, "tpa-my-topic-tsa", cfg.Kafka.Test.TopicName())
	})

	t.Run("producer/deprecated", func(t *testing.T) {
		var cfg struct {
			Kafka struct {
				Test configs.KafkaProducer `json:"test-4"`
			}
		}
		if !assert.NoError(t, provider.LoadConfig("kafka", &cfg)) {
			t.FailNow()
		}
		assert.Error(t, cfg.Kafka.Test.Bootstrap(provider))
	})
}
