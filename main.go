package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type ProcessState int

const (
	StateIdle ProcessState = iota
	StateDevRunning
	StateReviewRunning
	StateError
)

type Orchestrator struct {
	workDir     string
	coordDir    string
	watcher     *fsnotify.Watcher
	state       ProcessState
	stateMutex  sync.RWMutex
	currentProc *exec.Cmd
	procMutex   sync.Mutex
}

func NewOrchestrator(workDir string) (*Orchestrator, error) {
	coordDir := filepath.Join(workDir, ".claude-coordination")
	if err := os.MkdirAll(coordDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create coordination dir: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	orch := &Orchestrator{
		workDir:  workDir,
		coordDir: coordDir,
		watcher:  watcher,
		state:    StateIdle,
	}

	// Watch the coordination directory
	if err := watcher.Add(coordDir); err != nil {
		return nil, fmt.Errorf("failed to watch coordination dir: %w", err)
	}

	return orch, nil
}

func (o *Orchestrator) GetState() ProcessState {
	o.stateMutex.RLock()
	defer o.stateMutex.RUnlock()
	return o.state
}

func (o *Orchestrator) setState(state ProcessState) {
	o.stateMutex.Lock()
	defer o.stateMutex.Unlock()
	o.state = state
	log.Printf("State changed to: %v", state)
}

func (o *Orchestrator) Start() error {
	log.Printf("🚀 Starting orchestrator, watching: %s", o.coordDir)

	// Create initial trigger files if they don't exist
	o.initializeTriggerFiles()

	go o.watchFiles()
	return nil
}

func (o *Orchestrator) Stop() error {
	log.Printf("🛑 Stopping orchestrator...")

	o.procMutex.Lock()
	if o.currentProc != nil && o.currentProc.Process != nil {
		log.Printf("Terminating running process...")
		o.currentProc.Process.Kill()
	}
	o.procMutex.Unlock()

	return o.watcher.Close()
}

func (o *Orchestrator) initializeTriggerFiles() {
	triggers := []string{"task-ready.trigger", "dev-complete.trigger", "review-complete.trigger"}

	for _, trigger := range triggers {
		path := filepath.Join(o.coordDir, trigger)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Create empty trigger file
			if err := os.WriteFile(path, []byte(""), 0644); err != nil {
				log.Printf("Warning: failed to create %s: %v", trigger, err)
			}
		}
	}
}

func (o *Orchestrator) watchFiles() {
	for {
		select {
		case event, ok := <-o.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				o.handleFileChange(event.Name)
			}

		case err, ok := <-o.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("❌ Watcher error: %v", err)
		}
	}
}

func (o *Orchestrator) handleFileChange(filename string) {
	basename := filepath.Base(filename)

	log.Printf("📁 File changed: %s (current state: %v)", basename, o.GetState())

	switch basename {
	case "task-ready.trigger":
		if o.GetState() == StateIdle {
			o.startDevSession()
		}

	case "dev-complete.trigger":
		if o.GetState() == StateDevRunning {
			o.startReviewSession()
		}

	case "review-complete.trigger":
		if o.GetState() == StateReviewRunning {
			o.completeCycle()
		}
	}
}

func (o *Orchestrator) startDevSession() {
	log.Printf("🔨 Starting Developer session...")
	o.setState(StateDevRunning)

	// Run the Claude Code dev wrapper script
	devScript := filepath.Join(o.workDir, "claude-dev.sh")
	log.Printf("🔍 Dev script path: %s", devScript)
	o.startProcess(devScript, []string{o.workDir})
}

func (o *Orchestrator) startReviewSession() {
	log.Printf("👀 Starting Review session...")
	o.setState(StateReviewRunning)

	// Run the Claude Code review wrapper script
	reviewScript := filepath.Join(o.workDir, "claude-review.sh")
	log.Printf("🔍 Review script path: %s", reviewScript)
	o.startProcess(reviewScript, []string{o.workDir})
}

func (o *Orchestrator) completeCycle() {
	log.Printf("✅ Cycle complete, returning to idle")
	o.setState(StateIdle)
}

func (o *Orchestrator) startProcess(scriptPath string, args []string) {
	o.procMutex.Lock()
	defer o.procMutex.Unlock()

	// Debug: log what we're trying to run
	log.Printf("🔍 Attempting to run: %s with args: %v", scriptPath, args)
	log.Printf("🔍 Working directory: %s", o.workDir)

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		log.Printf("❌ Script not found at: %s", scriptPath)
		o.setState(StateError)
		return
	}

	// Start the process - run script directly
	cmd := exec.Command(scriptPath, args...)
	cmd.Dir = o.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Printf("❌ Failed to start %s: %v", scriptPath, err)
		o.setState(StateError)
		return
	}

	o.currentProc = cmd

	// Monitor process completion in background
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("❌ Process %s failed: %v", scriptPath, err)
			o.setState(StateError)
		} else {
			log.Printf("✅ Process %s completed successfully", scriptPath)
		}

		o.procMutex.Lock()
		o.currentProc = nil
		o.procMutex.Unlock()
	}()
}

// Simulate external triggers for testing
func (o *Orchestrator) TriggerTask() {
	triggerFile := filepath.Join(o.coordDir, "task-ready.trigger")
	o.touchFile(triggerFile)
}

func (o *Orchestrator) TriggerDevComplete() {
	triggerFile := filepath.Join(o.coordDir, "dev-complete.trigger")
	o.touchFile(triggerFile)
}

func (o *Orchestrator) TriggerReviewComplete() {
	triggerFile := filepath.Join(o.coordDir, "review-complete.trigger")
	o.touchFile(triggerFile)
}

func (o *Orchestrator) touchFile(filename string) {
	// Write current timestamp to trigger file change event
	timestamp := fmt.Sprintf("triggered at %s\n", time.Now().Format(time.RFC3339))
	if err := os.WriteFile(filename, []byte(timestamp), 0644); err != nil {
		log.Printf("Failed to touch %s: %v", filename, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <work-directory> [--demo]")
	}

	workDir := os.Args[1]
	demoMode := len(os.Args) > 2 && os.Args[2] == "--demo"

	orch, err := NewOrchestrator(workDir)
	if err != nil {
		log.Fatal("Failed to create orchestrator:", err)
	}
	defer orch.Stop()

	if err := orch.Start(); err != nil {
		log.Fatal("Failed to start orchestrator:", err)
	}

	if demoMode {
		// Demo: simulate the workflow
		log.Printf("🎮 Starting demo workflow...")

		time.Sleep(1 * time.Second)
		log.Printf("📝 Triggering new task...")
		orch.TriggerTask()

		time.Sleep(3 * time.Second)
		log.Printf("✅ Simulating dev completion...")
		orch.TriggerDevComplete()

		time.Sleep(3 * time.Second)
		log.Printf("✅ Simulating review completion...")
		orch.TriggerReviewComplete()
	} else {
		log.Printf("📝 To start a task: echo 'start' > .claude-coordination/task-ready.trigger")
	}

	// Keep running to observe
	log.Printf("👂 Watching for file changes... (press Ctrl+C to exit)")
	select {} // Block forever
}
