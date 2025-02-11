package configs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vredens/infrastructure"
	"github.com/vredens/infrastructure/configs"
)

func TestNSQConsumer(t *testing.T) {
	provider, err := infrastructure.NewProvider(infrastructure.ProviderSettings{
		EnvName:       "nsq-tests",
		SystemName:    "sys",
		ComponentName: "comp",
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Run("consumer-1", func(t *testing.T) {
		var config struct {
			Test configs.NSQConsumer `json:"valid-consumer-1"`
		}
		if !assert.NoError(t, provider.LoadConfig("nsq", &config)) {
			t.FailNow()
		}
		if !assert.NoError(t, config.Test.Bootstrap(provider)) {
			t.FailNow()
		}
		assert.Equal(t, "nsq-tests-sys-topic-1-comp", config.Test.TopicName())
		assert.Equal(t, "nsq-tests-sys-channel-1-comp", config.Test.ChannelName())
	})

	t.Run("consumer-2", func(t *testing.T) {
		var config struct {
			Test configs.NSQConsumer `json:"valid-consumer-2"`
		}
		if !assert.NoError(t, provider.LoadConfig("nsq", &config)) {
			t.FailNow()
		}
		if !assert.NoError(t, config.Test.Bootstrap(provider)) {
			t.FailNow()
		}
		assert.Equal(t, "topic-2", config.Test.TopicName())
		assert.Equal(t, "channel-2", config.Test.ChannelName())
	})

	t.Run("producer-1", func(t *testing.T) {
		var config struct {
			Test configs.NSQProducer `json:"valid-producer-1"`
		}
		if !assert.NoError(t, provider.LoadConfig("nsq", &config)) {
			t.FailNow()
		}
		if !assert.NoError(t, config.Test.Bootstrap(provider)) {
			t.FailNow()
		}
		assert.Equal(t, "nsq-tests-sys-topic-1-comp", config.Test.TopicName())
	})
}
