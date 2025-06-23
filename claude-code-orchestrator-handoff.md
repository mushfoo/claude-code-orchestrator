# Claude Code Handoff: Two-Instance Collaboration Orchestrator

## Project Overview

Building a Go-based orchestrator that coordinates two Claude Code instances in a peer review workflow:

- **Instance A (Developer)**: Implements features using project guidelines
- **Instance B (Reviewer)**: Reviews code against quality principles, catches over-engineering

## Technical Validation Results ✅

We built and thoroughly tested a file-watching orchestration system. **All validation tests passed:**

### Core Architecture Validated

- **File-based coordination**: Robust communication through trigger files
- **Process lifecycle management**: Clean spawn, monitor, and cleanup of processes
- **State machine**: Reliable transitions `Idle → DevRunning → ReviewRunning → Idle`
- **Error handling**: Graceful failure detection and recovery
- **Performance**: Sub-second response to file changes, no race conditions

### Stress Testing Results

1. **Happy Path**: Perfect state transitions and file coordination
2. **Manual Control**: Out-of-order triggers handled correctly
3. **Race Conditions**: No issues with rapid fire triggers (0.2s intervals)
4. **Error Recovery**: Process failures detected immediately with proper state management
5. **Process Interruption**: External kills handled gracefully

## Technical Implementation

### Proven Architecture

```
Project Directory/
├── src/                           # Actual project code
├── .claude-coordination/          # Coordination layer
│   ├── current-task.md           # From toolkit task breakdown
│   ├── dev-session.md            # Developer's work log
│   ├── review-request.md         # Code ready for review
│   ├── review-feedback.md        # Reviewer's comments
│   └── project-context.md        # PRD + guidelines (shared)
├── templates/
│   ├── dev-prompt.md             # Developer session template
│   └── review-prompt.md          # Reviewer session template
└── orchestrator                  # Go binary
```

### Validated Go Components

- **File watcher**: Uses `github.com/fsnotify/fsnotify` - handles rapid changes reliably
- **Process management**: `exec.Command` with proper cleanup - no zombie processes
- **State synchronization**: Mutex-protected state machine - no race conditions
- **Error handling**: Exit code detection with graceful workflow halting

### Working Implementation Pattern

```go
// Trigger files for coordination
task-ready.trigger    → Starts Developer instance
dev-complete.trigger  → Starts Reviewer instance
review-complete.trigger → Returns to idle

// Process spawning (macOS compatible)
cmd := exec.Command(scriptPath, args...)
cmd.Dir = workDir
cmd.Start()
// Background monitoring with proper cleanup
```

## Next Implementation Phase

### 1. Claude Code Integration (High Priority)

Replace simulation scripts with actual Claude Code calls:

```bash
# Replace simulate-dev.sh with:
claude-code --prompt-file templates/dev-prompt.md --context-file .claude-coordination/current-task.md

# Replace simulate-review.sh with:
claude-code --prompt-file templates/review-prompt.md --context-file .claude-coordination/dev-session.md
```

### 2. Prompt Templates (Based on Toolkit Principles)

Create role-specific prompts that prevent over-engineering:

**Dev Prompt Template**:

- Reference project PRD and current task
- Emphasize simplest solution that works
- Include testing requirements
- Output structured session log

**Review Prompt Template**:

- Reference development session output
- Check against over-engineering patterns
- Validate real-world usability
- Output structured feedback

### 3. Enhanced Error Handling

- **Process timeouts**: Claude Code sessions shouldn't run indefinitely
- **Retry logic**: Handle temporary Claude Code failures
- **Conflict resolution**: What happens when reviewer rejects repeatedly?

### 4. Configuration Management

- **Configurable paths**: Claude Code binary location, templates, etc.
- **Timeout settings**: Per-session and overall workflow timeouts
- **Logging levels**: Debug, info, error output control

## Implementation Requirements

### Must Haves

1. **Absolute path handling**: Relative paths cause issues (validated)
2. **Process cleanup**: Prevent zombie processes (validated pattern works)
3. **Error state recovery**: System must recover from any failure state
4. **File locking**: Prevent coordination file corruption during concurrent writes

### Quality Standards

- **No over-engineering**: Follow the toolkit's principle of simplest solution
- **Test coverage**: Unit tests for state machine, integration tests for full workflow
- **Real-world validation**: Actually use the system to develop a small project
- **Documentation**: Clear setup and usage instructions

## Proven Development Approach

Based on validation testing:

1. **Start simple**: Get basic Claude Code integration working first
2. **Incremental complexity**: Add timeouts, retries, configuration iteratively
3. **Test continuously**: Validate each addition with real scenarios
4. **Use absolute paths**: `$(pwd)` or full path resolution for all file operations

## File Handoff Structure

### Context Files to Create

```markdown
# .claude-coordination/project-context.md

- Project PRD summary
- Development guidelines from toolkit
- Key architectural decisions

# .claude-coordination/current-task.md

- Specific task from toolkit breakdown
- Acceptance criteria
- Timeline and priority

# templates/dev-prompt.md

- Role definition for Developer instance
- Reference to project context and current task
- Output format requirements

# templates/review-prompt.md

- Role definition for Reviewer instance
- Quality checklist based on toolkit principles
- Feedback format requirements
```

## Success Metrics

### Technical

- **Zero race conditions**: File coordination remains reliable under stress
- **Clean error recovery**: System returns to working state after any failure
- **Performance**: < 1 second response to triggers, < 30 seconds per Claude Code session

### Workflow

- **Over-engineering prevention**: Reviewer catches unnecessarily complex solutions
- **Quality improvement**: Code quality measurably improves through review iterations
- **Context preservation**: Full project context maintained across sessions

## Ready for Development

The orchestration foundation is **validated and robust**. Focus can now shift entirely to Claude Code integration and prompt engineering, knowing the coordination mechanics work reliably.

**Recommended first commit**: Replace simulation scripts with basic Claude Code calls using validated file coordination pattern.
