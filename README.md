# Infrastructure configuration provider

This package aims to be a central point for configurations and helpers around infrastructure resources.

There are two types of configurations: infrastructure and application. Best way to know if something belongs to the infrastructure or not is to go through the checklist

- I need these configurations to access resources external to my program such as SaaS services, databases, webservices, etc.
- Almost every time I need to deploy to a new environment I need a new value for these configurations

If either answer to the above was true then you are most certainly looking at an infrastructure configuration.

Application configurations, on the other hand, are usually immutable across environments.

**Infrastructure configuration**

Infrastructure resources have a specific configuration associated with the environment. These resources are identified across environments by their ARN (Application Resource Name). An ARN has the following format: `arn://<group>[/<subgroup>]/<slug>`. Example: `arn://storage/postgres/users`.

The infrastructure resource data structures reside in [./resources](/resources)`. There is a full example of a JSON file with the structure for configuring these settings at [testdata/infra/test.json](./testdata/infra/test.json).

These settings should remain immutable through time, unlike application configurations.

**Application configurations**

Application configurations can make reference to a ARN in order to know where to locate a certain resource. This means that application configurations should remain immutable between environments unless specific tweaking is necessary.

There are examples of configurations for each resource under folder [./configs](/configs). You can use your own version of these in your application.

There's an example of an application configuration file at [testdata/app.json](./testdata/config/app.json).

You can add a specific application configuration for a certain environment. For example, if you have a `my-app.json` configuration file you can create a custom configuration for the `dev` environment by creating a copy of that configuration and naming it `my-app.dev.json`. This will **not** mix in configurations

**Certificate Authorities**

You can add custom CA certificates to the system wide list of CAs which can then be used to configure HTTP connections. The first location where a valid certificate is found is the only location used.

The default locations to look for custom certificates are, in order:

- `/etc/certs`
- `./etc/certs`
- `./testdata/certs`

Also, files names must match

- `ca.pem`
- `*.ca.pem`

The main reason is that in the future we might add support for loading other types of certificates (e.g.: server/client).

### Secrets

The infrastructure package support render values from environment variables. In order to support multiple values and have some flexibility we leverage templates to render secrets when loading the configuration.

The structure used when rendering the infrastructure configuration file is: 

```
{
	Environment: "",
	System: "",
	Component: "",
	Env: {
		...
	}
}
```

* `Env` contains all environment variables.
* `Environment` is set when configuring the Provider.
* `System` is set when configuring the Provider.
* `Component` is set when configuring the Provider.

## Quick Start

```golang
package main

import (
	"fmt"

	"github.com/vredens/infrastructure"
	"github.com/vredens/infrastructure/resources"
)

type MyAppConfig struct {
	Elastic resources.ElasticsearchConfig `json:"elastic"`
}

func main() {
	provider, err := infrastructure.NewProvider()
	if err != nil {
		panic(err)
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
```
