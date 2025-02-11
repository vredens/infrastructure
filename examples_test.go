package infrastructure_test

import (
	"fmt"
	"os"

	"github.com/vredens/infrastructure"
	"github.com/vredens/infrastructure/configs"
)

func ExampleNewProvider() {
	provider, err := infrastructure.NewProvider(infrastructure.ProviderSettings{
		EnvName:       infrastructure.GetFromEnv("SYSTEM_ENV", "ENV"),
		SystemName:    os.Getenv("SYSTEM"),
		ComponentName: os.Getenv("COMPONENT"),
	})
	if err != nil {
		panic(err)
	}

	type MyAppConfig struct {
		Elastic configs.Elasticsearch `json:"elastic"`
	}

	// load a configuration file
	var config MyAppConfig
	if err := provider.LoadConfig("myapp", &config); err != nil {
		panic(err)
	}

	// bootstrap the configuration which loads the infra resource
	if err := config.Elastic.Bootstrap(provider); err != nil {
		panic(err)
	}
	fmt.Printf("create a conneciton using: %+v\n", config.Elastic.Resource())
	fmt.Printf("the config has extra params: %+v\n", config.Elastic.Params)
}

func ExampleProvider_Locator() {
	provider, err := infrastructure.NewProvider(infrastructure.ProviderSettings{
		EnvName:       infrastructure.GetFromEnv("APP_ENV", "ENV"),
		SystemName:    os.Getenv("SYSTEM_NAME"),
		ComponentName: os.Getenv("COMPONENT_NAME"),
	})
	if err != nil {
		panic(err)
	}

	// a custom configuration structure that has a ARN.
	type MyAppConfig struct {
		ResourceName string `json:"arn"`
		Index        string `json:"index"`
	}

	// load a configuration file
	var config MyAppConfig
	if err := provider.LoadConfig("myapp", &config); err != nil {
		panic(err)
	}

	// fetching a resource directly by its ARN.
	resource := provider.Locator().LocateElasticResource(config.ResourceName)
	if err := resource.Error(); err != nil {
		panic(err)
	}
	fmt.Printf("create a conneciton using: %+v\n", resource)
}
