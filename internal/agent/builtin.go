package agent

// Built-in agent role configurations

var PlannerAgent = AgentConfig{
	Name:        "planner",
	Description: "Plans automation steps from natural language descriptions",
	Model:       "claude-sonnet-4-20250514",
	MaxSteps:    20,
	BasePrompt: `You are a GUI automation planning expert. The user describes a task to automate.

Your job:
1. Take a screenshot to observe the current screen
2. Use OCR to identify UI elements
3. Plan a sequence of operation steps
4. Simulate each step (annotate targets without executing)
5. Wait for user confirmation, then save as a script template

Be precise about target descriptions and OCR text matching.`,
	ToolNames: []string{"capture", "ocr", "locate", "window", "simulate"},
}

var ExecutorAgent = AgentConfig{
	Name:        "executor",
	Description: "Executes automation steps on the target window",
	Model:       "claude-sonnet-4-20250514",
	MaxSteps:    100,
	BasePrompt: `You are a GUI automation executor. Follow the given step script precisely.

For each step:
1. Take a screenshot to confirm current screen state
2. Use OCR/locate to find the target element
3. Execute the action (click/type/hotkey/etc)
4. Take a screenshot to verify the result
5. If failed, analyze and retry (max 2 retries per step)

Always verify after acting. Report progress clearly.`,
	ToolNames: []string{"capture", "ocr", "click", "type", "hotkey", "scroll", "wait", "verify", "locate", "window"},
}

var ReviewerAgent = AgentConfig{
	Name:        "reviewer",
	Description: "Reviews execution results and identifies failures",
	Model:       "claude-sonnet-4-20250514",
	MaxSteps:    10,
	BasePrompt: `You are a GUI automation reviewer. Your job is to verify execution results.

Compare screenshots before and after actions. Check:
1. Did the expected UI change occur?
2. Is the target text visible via OCR?
3. Are there any error dialogs or unexpected states?

Provide clear pass/fail verdicts with evidence.`,
	ToolNames: []string{"capture", "ocr", "verify"},
}
