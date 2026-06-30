package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/yc446833448/VisuTask/internal/action"
	"github.com/yc446833448/VisuTask/internal/agent"
	"github.com/yc446833448/VisuTask/internal/agent/tools"
	"github.com/yc446833448/VisuTask/internal/config"
	"github.com/yc446833448/VisuTask/internal/concurrency"
	"github.com/yc446833448/VisuTask/internal/llm"
	"github.com/yc446833448/VisuTask/internal/monitor"
	"github.com/yc446833448/VisuTask/internal/store"
	"github.com/yc446833448/VisuTask/internal/vision"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database
	db, err := store.New(cfg.GetDBPath())
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Initialize vision engine
	screenCapturer := vision.NewScreenCapturer()

	var ocrRecognizer vision.Recognizer
	remoteOCR, err := vision.NewRemoteOCRRecognizerWithConfig(cfg.OCR.Endpoint, cfg.OCR.APIKey)
	if err != nil {
		log.Printf("warning: Remote OCR not available (%v), using stub", err)
		ocrRecognizer = vision.NewStubRecognizer()
	} else {
		ocrRecognizer = remoteOCR
		log.Printf("OCR: using remote service at %s", cfg.OCR.Endpoint)
	}

	visionEngine := vision.NewEngine(screenCapturer, ocrRecognizer, vision.NewStubDetector())

	// Initialize action engine (platform-specific: Windows=Win32API, macOS=CoreGraphics+AppleScript, other=stub)
	actionEngine := action.NewEngine(
		action.NewPlatformMouse(),
		action.NewPlatformKeyboard(),
		action.NewPlatformWindow(),
	)

	// On macOS, check and request Accessibility permissions
	action.EnsureAccessibility()

	monitorChecker := monitor.NewChecker(visionEngine)

	// Initialize LLM Gateway with configured providers
	gateway := initLLMGateway(cfg)

	// Initialize concurrency manager
	user, err := db.GetUser()
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	concurrencyMgr := concurrency.NewManager(user.GetMaxConcurrent())

	// Initialize tool registry
	registry := agent.NewToolRegistry()
	tools.RegisterAll(registry, visionEngine, actionEngine, monitorChecker)

	// Setup agent tool assignments
	registry.AssignAgentTools("planner", agent.PlannerAgent.ToolNames)
	registry.AssignAgentTools("executor", agent.ExecutorAgent.ToolNames)
	registry.AssignAgentTools("reviewer", agent.ReviewerAgent.ToolNames)

	// Create agent service
	agentService := agent.NewService(db, visionEngine, actionEngine, monitorChecker, gateway, concurrencyMgr, registry)

	// Create Wails application
	app := application.New(application.Options{
		Name:        "VisuTask",
		Description: "GUI Automation Desktop App",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Services: []application.Service{
			application.NewService(agentService),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	// Create main window (must use app.Window.NewWithOptions, NOT application.NewWindow)
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "VisuTask",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
	})

	// Bridge agent events → Wails events
	agentService.Events().On(func(event agent.EventType, data interface{}) {
		app.Event.Emit(string(event), data)
	})

	log.Println("VisuTask starting...")
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// initLLMGateway creates the LLM gateway based on config
func initLLMGateway(cfg *config.Config) *llm.Gateway {
	var providers []llm.Provider

	for _, name := range cfg.LLM.FallbackOrder {
		switch name {
		case "anthropic":
			if cfg.LLM.Anthropic.APIKey != "" {
				p := llm.NewAnthropicProvider()
				p.SetConfig(cfg.LLM.Anthropic.APIKey, cfg.LLM.Anthropic.BaseURL)
				providers = append(providers, p)
				log.Printf("LLM: Anthropic enabled (model: %s)", cfg.LLM.Anthropic.Model)
			}
		case "openai":
			if cfg.LLM.OpenAI.APIKey != "" {
				p := llm.NewOpenAIAdapter()
				p.SetConfig(cfg.LLM.OpenAI.APIKey, cfg.LLM.OpenAI.BaseURL, "openai")
				providers = append(providers, p)
				log.Printf("LLM: OpenAI enabled (model: %s)", cfg.LLM.OpenAI.Model)
			}
		case "ollama":
			p := llm.NewOllamaAdapter()
			p.SetOllamaConfig(cfg.LLM.Ollama.BaseURL)
			providers = append(providers, p)
			log.Printf("LLM: Ollama enabled (model: %s, url: %s)", cfg.LLM.Ollama.Model, cfg.LLM.Ollama.BaseURL)
		}
	}

	if len(providers) == 0 {
		log.Println("warning: no LLM providers configured, set API keys in config.yaml")
	}

	return llm.NewGateway(providers...)
}
