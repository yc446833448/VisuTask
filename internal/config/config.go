package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the top-level application configuration
type Config struct {
	LLM         LLMConfig         `yaml:"llm"`
	OCR         OCRConfig         `yaml:"ocr"`
	Agent       AgentConfig       `yaml:"agent"`
	Retry       RetryConfig       `yaml:"retry"`
	Concurrency ConcurrencyConfig `yaml:"concurrency"`
	Scripts     ScriptsConfig     `yaml:"scripts"`
	Database    DatabaseConfig    `yaml:"database"`
}

type LLMConfig struct {
	Primary       string          `yaml:"primary"`
	Anthropic     ProviderConfig  `yaml:"anthropic"`
	OpenAI        ProviderConfig  `yaml:"openai"`
	Ollama        ProviderConfig  `yaml:"ollama"`
	FallbackOrder []string        `yaml:"fallback_order"`
}

type ProviderConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
	Model   string `yaml:"model"`
}

type OCRConfig struct {
	Endpoint string `yaml:"endpoint"`
	APIKey   string `yaml:"api_key"`
}

type AgentConfig struct {
	PlannerMaxSteps   int `yaml:"planner_max_steps"`
	ExecutorMaxSteps  int `yaml:"executor_max_steps"`
	ReviewerMaxSteps  int `yaml:"reviewer_max_steps"`
	ToolTimeout       int `yaml:"tool_timeout"`
	DoomLoopThreshold int `yaml:"doom_loop_threshold"`
	ContextWindow     int `yaml:"context_window"`
}

type RetryConfig struct {
	InitialDelayMS int     `yaml:"initial_delay_ms"`
	MaxDelayMS     int     `yaml:"max_delay_ms"`
	BackoffFactor  float64 `yaml:"backoff_factor"`
	MaxRetries     int     `yaml:"max_retries"`
}

type ConcurrencyConfig struct {
	DefaultMax int `yaml:"default_max"`
}

type ScriptsConfig struct {
	Directory string `yaml:"directory"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// DefaultConfig returns the configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			Primary:       "anthropic",
			FallbackOrder: []string{"anthropic", "openai", "ollama"},
			Anthropic: ProviderConfig{
				Model: "claude-sonnet-4-20250514",
			},
			OpenAI: ProviderConfig{
				Model: "gpt-4o",
			},
			Ollama: ProviderConfig{
				BaseURL: "http://localhost:11434",
				Model:   "llama3",
			},
		},
		Agent: AgentConfig{
			PlannerMaxSteps:   20,
			ExecutorMaxSteps:  100,
			ReviewerMaxSteps:  10,
			ToolTimeout:       30,
			DoomLoopThreshold: 3,
			ContextWindow:     128000,
		},
		Retry: RetryConfig{
			InitialDelayMS: 2000,
			MaxDelayMS:     30000,
			BackoffFactor:  2.0,
			MaxRetries:     5,
		},
		Concurrency: ConcurrencyConfig{
			DefaultMax: 3,
		},
	}
}

// Load reads the configuration from the YAML file.
// If the file doesn't exist, returns default config.
// Environment variables override file values.
func Load() (*Config, error) {
	cfg := DefaultConfig()

	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No config file, use defaults + env overrides
			applyEnvOverrides(cfg)
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Fill in defaults for zero values
	applyDefaults(cfg)

	// Environment variables take priority
	applyEnvOverrides(cfg)

	return cfg, nil
}

// getConfigPath returns the config file path
func getConfigPath() string {
	// Check VISUTASK_CONFIG env first
	if p := os.Getenv("VISUTASK_CONFIG"); p != "" {
		return p
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "config.yaml"
	}
	return filepath.Join(home, ".visutask", "config.yaml")
}

// applyDefaults fills zero values with sensible defaults
func applyDefaults(cfg *Config) {
	if cfg.Agent.PlannerMaxSteps == 0 {
		cfg.Agent.PlannerMaxSteps = 20
	}
	if cfg.Agent.ExecutorMaxSteps == 0 {
		cfg.Agent.ExecutorMaxSteps = 100
	}
	if cfg.Agent.ReviewerMaxSteps == 0 {
		cfg.Agent.ReviewerMaxSteps = 10
	}
	if cfg.Agent.ToolTimeout == 0 {
		cfg.Agent.ToolTimeout = 30
	}
	if cfg.Agent.DoomLoopThreshold == 0 {
		cfg.Agent.DoomLoopThreshold = 3
	}
	if cfg.Agent.ContextWindow == 0 {
		cfg.Agent.ContextWindow = 128000
	}
	if cfg.Retry.MaxRetries == 0 {
		cfg.Retry.MaxRetries = 5
	}
	if cfg.Retry.BackoffFactor == 0 {
		cfg.Retry.BackoffFactor = 2.0
	}
	if cfg.Concurrency.DefaultMax == 0 {
		cfg.Concurrency.DefaultMax = 3
	}
	if len(cfg.LLM.FallbackOrder) == 0 {
		cfg.LLM.FallbackOrder = []string{"anthropic", "openai", "ollama"}
	}
}

// applyEnvOverrides lets environment variables override config values
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("ANTHROPIC_API_KEY"); v != "" {
		cfg.LLM.Anthropic.APIKey = v
	}
	if v := os.Getenv("ANTHROPIC_BASE_URL"); v != "" {
		cfg.LLM.Anthropic.BaseURL = v
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.LLM.OpenAI.APIKey = v
	}
	if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
		cfg.LLM.OpenAI.BaseURL = v
	}
	if v := os.Getenv("OLLAMA_BASE_URL"); v != "" {
		cfg.LLM.Ollama.BaseURL = v
	}
	if v := os.Getenv("OCR_ENDPOINT"); v != "" {
		cfg.OCR.Endpoint = v
	}
	if v := os.Getenv("OCR_API_KEY"); v != "" {
		cfg.OCR.APIKey = v
	}
}

// GetDBPath returns the resolved database path
func (c *Config) GetDBPath() string {
	if c.Database.Path != "" {
		return c.Database.Path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "visutask.db"
	}
	dir := filepath.Join(home, ".visutask")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "visutask.db")
}

// GetScriptsDir returns the resolved scripts directory
func (c *Config) GetScriptsDir() string {
	if c.Scripts.Directory != "" {
		return c.Scripts.Directory
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "./scripts"
	}
	dir := filepath.Join(home, ".visutask", "scripts")
	os.MkdirAll(dir, 0755)
	return dir
}
