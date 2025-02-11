package infrastructure

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/vredens/infrastructure/lib/certs"
	"github.com/vredens/infrastructure/resources"
)

var defaults ProviderSettings = ProviderSettings{
	CertFolders: []string{
		"/etc/certs",
		"etc/certs",
		"testdata/certs",
	},
	AppConfigFolders: []string{
		"/etc/config",
		"etc/config",
		"testdata/config",
	},
	InfraConfigFolders: []string{
		"/etc/infra",
		"etc/infra",
		"testdata/infra",
	},
	EnvName:       "local",
	SystemName:    "NOSYSTEM",
	ComponentName: "NOCOMPONENT",
}

// ProviderSettings allows setting some parameters used by the infrastructure provider to load configurations.
type ProviderSettings struct {
	// EnvName is used to select configuration files. Defaults to `local`.
	EnvName string
	// SystemName is the name by which your system is known. Defaults to `NOSYSTEM`.
	SystemName string
	// ComponentName is the name by which the running process is known, which can be part of a larger system.
	// Defaults to `NOCOMPONENT`.
	ComponentName string
	// EnvVarPrefix is the prefix used to look for environment variables.
	// If empty then all environment variables are used.
	EnvVarPrefix string
	// CertFolders are locations where to look for certificate files (*.pem, etc).
	// These folders are used by provider.Certs().
	// Defaults to /etc/certs, etc/certs, testdata/certs.
	CertFolders []string
	// AppConfigFolders are locations where to look for application configuration files.
	// These locations are used by methods such as provider.LoadConfig().
	// Defaults to /etc/config, etc/config, testdata/config.
	AppConfigFolders []string
	// InfraConfigFolders are locations where to look for infrastructure configuration files.
	// These locations are used when initializing a new Provider.
	// Defaults to /etc/infra, etc/infra, testdata/infra.
	InfraConfigFolders []string
}

func (settings ProviderSettings) sanitize() ProviderSettings {
	if len(settings.AppConfigFolders) == 0 {
		settings.AppConfigFolders = defaults.AppConfigFolders
	}
	if len(settings.InfraConfigFolders) == 0 {
		settings.InfraConfigFolders = defaults.InfraConfigFolders
	}
	if len(settings.CertFolders) == 0 {
		settings.CertFolders = defaults.CertFolders
	}
	if settings.EnvName == "" {
		settings.EnvName = defaults.EnvName
	}
	if settings.SystemName == "" {
		settings.SystemName = defaults.SystemName
	}
	if settings.ComponentName == "" {
		settings.ComponentName = defaults.ComponentName
	}
	return settings
}

func (settings ProviderSettings) Validate() error {
	if settings.EnvName == "" {
		return errors.New("could not determine environment name")
	}
	if settings.SystemName == "" {
		return errors.New("could not determine system name from environment")
	}
	if settings.ComponentName == "" {
		return errors.New("could not determine component name from environment")
	}
	return nil
}

// Provider of infrastructure resources and repo of application settings.
type Provider struct {
	infraConfig  resources.Locator
	cfgLoader    *viper.Viper
	settings     ProviderSettings
	resourcePath string
	certs        certs.Certs
	data         tmplData
}

type tmplData struct {
	Environment string
	System      string
	Component   string
	Env         map[string]string
}

// NewProvider creates a new infrastructure provider with the given settings.
func NewProvider(settings ProviderSettings) (*Provider, error) {
	provider := &Provider{
		settings: settings.sanitize(),
	}
	if err := provider.settings.Validate(); err != nil {
		return nil, fmt.Errorf("invalid environment settings; %w", err)
	}
	provider.data = tmplData{
		Environment: settings.EnvName,
		System:      settings.SystemName,
		Component:   settings.ComponentName,
		Env:         make(map[string]string),
	}

	for _, envVar := range os.Environ() {
		if settings.EnvVarPrefix == "" || strings.HasPrefix(envVar, settings.EnvVarPrefix) {
			parts := strings.SplitN(envVar, "=", 2)
			provider.data.Env[parts[0]] = parts[1]
		}
	}

	provider.cfgLoader = viper.New()
	provider.cfgLoader.SetConfigType("json")
	for _, path := range provider.settings.AppConfigFolders {
		provider.cfgLoader.AddConfigPath(path)
	}

	vInfra, err := provider.loadInfraConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create infra; %w", err)
	}
	provider.resourcePath = vInfra.ConfigFileUsed()
	provider.infraConfig.SetProvider(provider)
	if err := vInfra.Unmarshal(&provider.infraConfig, func(cfg *mapstructure.DecoderConfig) { cfg.TagName = "json" }); err != nil {
		return provider, fmt.Errorf("failed to unmarshal infrastructure configuration; %w", err)
	}

	provider.certs = certs.New(certs.Config{Locations: provider.settings.CertFolders})

	return provider, nil
}

func (provider *Provider) loadInfraConfig() (empty *viper.Viper, err error) {
	vInfra := viper.New()
	vInfra.SetConfigName(provider.settings.EnvName)
	vInfra.SetConfigType("json")
	for _, path := range provider.settings.InfraConfigFolders {
		vInfra.AddConfigPath(path)
	}
	if err := vInfra.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in infrastructure resource configuration; %w", err)
	}

	fileContent, err := os.ReadFile(vInfra.ConfigFileUsed())
	if err != nil {
		return empty, fmt.Errorf("failed to read config file %s", vInfra.ConfigFileUsed())
	}
	renderedConfig, err := provider.RenderSecrets(string(fileContent))
	if err != nil {
		return empty, fmt.Errorf("failed to render secrets; %w", err)
	}

	err = vInfra.ReadConfig(strings.NewReader(renderedConfig))
	if err != nil {
		return empty, fmt.Errorf("failed to read config from rendered config; %w", err)
	}

	return vInfra, nil
}

// RenderSecrets using provider replaceVariables in the given config file
// it will return an error when the config file has invalid placeholerds, such as
// nonexisting functions (e.g. {{ test }}) or invalid properties (e.g. {{ .test }})
func (provider Provider) RenderSecrets(value string) (empty string, err error) {
	t, err := template.New("secrets").Parse(value)
	if err != nil {
		return empty, fmt.Errorf("failed create template from config file; %w", err)
	}

	tmp, err := renderTemplate(t, provider.data)
	if err != nil {
		return empty, fmt.Errorf("failed render template; %w", err)
	}
	return string(tmp), nil
}

// RenderSecrets using provider replaceVariables in the given config file returning
// given value if doesn't exist or is an invalid template function
func (provider Provider) RenderSecret(value string) string {
	t, err := template.New("secrets").Parse(value)
	if err != nil {
		return value
	}

	tmp, err := renderTemplate(t, provider.data)
	if err != nil {
		return value
	}
	return string(tmp)
}

func renderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	// TODO: this should be a buffer pool and 2048 is based on the amazing science of wild guess.
	buf := bytes.NewBuffer(make([]byte, 0, 2048))

	if err := tmpl.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template; %w", err)
	}

	return buf.Bytes(), nil
}

// LoadConfigFromTemplate into the config structure provided.
func (provider *Provider) LoadConfigFromTemplate(template []byte, config interface{}) error {
	renderedConfig, err := provider.RenderSecrets(string(template))
	if err != nil {
		return fmt.Errorf("failed to render secrets; %w", err)
	}

	err = provider.cfgLoader.ReadConfig(strings.NewReader(renderedConfig))
	if err != nil {
		return fmt.Errorf("failed to read rendered configuration")
	}

	if err := provider.cfgLoader.Unmarshal(config, func(cfg *mapstructure.DecoderConfig) { cfg.TagName = "json" }); err != nil {
		return fmt.Errorf("fail to unmarshal json configuration into the provided structure; %w", err)
	}
	return nil
}

func isConfigNotFound(err error) bool {
	var expected viper.ConfigFileNotFoundError
	return errors.As(err, &expected)
}

// LoadConfig into the config structure provided.
func (provider *Provider) LoadConfig(namespace string, config interface{}) error {
	global, err := provider.loadConfig(namespace, config)
	if err != nil {
		return fmt.Errorf("failed to load global configuration file; %w", err)
	}

	env, err := provider.loadConfig(namespace+"."+provider.settings.EnvName, config)
	if err != nil {
		return fmt.Errorf("failed to load environment configuration file; %w", err)
	}

	if !global && !env {
		return fmt.Errorf("no configuration found")
	}
	return nil
}

func (provider *Provider) loadConfig(name string, config interface{}) (loaded bool, err error) {
	provider.cfgLoader.SetConfigName(name)
	if err := provider.cfgLoader.ReadInConfig(); err != nil {
		if isConfigNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("read failed; %w", err)
	}
	if err := provider.LoadConfigFromFile(provider.cfgLoader.ConfigFileUsed(), config); err != nil {
		return false, fmt.Errorf("load failed; %w", err)
	}
	return true, nil
}

// LoadConfigFromFile into the config structure provided.
func (provider *Provider) LoadConfigFromFile(path string, config interface{}) error {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s", path)
	}
	renderedConfig, err := provider.RenderSecrets(string(fileContent))
	if err != nil {
		return fmt.Errorf("failed to render secrets; %w", err)
	}

	err = provider.cfgLoader.ReadConfig(strings.NewReader(renderedConfig))
	if err != nil {
		return fmt.Errorf("failed to read rendered configuration")
	}

	if err := provider.cfgLoader.Unmarshal(config, func(cfg *mapstructure.DecoderConfig) { cfg.TagName = "json" }); err != nil {
		return fmt.Errorf("fail to unmarshal json configuration into the provided structure; %w", err)
	}
	return nil
}

func (provider Provider) Certs() certs.Certs {
	return provider.certs
}

// ResourcePath of the current configuration
func (provider *Provider) ResourcePath() string {
	return provider.resourcePath
}

// Locator gives access to the infrastructure configuration for implementing your own providers.
func (provider *Provider) Locator() *resources.Locator {
	return &provider.infraConfig
}

// SystemName your process is running in.
func (provider *Provider) SystemName() string {
	return provider.settings.SystemName
}

// ComponentName your process is running in.
func (provider *Provider) ComponentName() string {
	return provider.settings.ComponentName
}

// Environment your process is running in.
func (provider *Provider) Environment() string {
	return provider.settings.EnvName
}

// RegisterCallback for notification when a configuration changes.
func (provider *Provider) RegisterCallback(fn func()) {

}
