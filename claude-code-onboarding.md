# Claude Code Onboarding: Two-Instance Orchestrator Implementation

## Project Context

You are continuing development of a **Go-based orchestrator** that coordinates two Claude Code instances in a peer review workflow for software development. The core orchestration mechanics have been **validated through comprehensive testing**.

### Validated Foundation

- ✅ File-watching coordination system (fsnotify)
- ✅ Process lifecycle management (exec.Command)
- ✅ State machine with mutex synchronization
- ✅ Error handling and recovery
- ✅ Race condition resistance under stress testing

### Architecture Overview

```
Orchestrator monitors trigger files → Spawns Claude Code instances → Instances coordinate via files
```

## Current Implementation Status

### ✅ COMPLETED (Do Not Rebuild)

**File Watcher Core** - Fully functional and tested:

- State machine: `Idle → DevRunning → ReviewRunning → Idle`
- File monitoring with `github.com/fsnotify/fsnotify`
- Process spawning and cleanup
- Error detection and state recovery
- Comprehensive validation across 5 test scenarios

**Working Go Code Structure** - Location: `spike-test/main.go`

- Orchestrator struct with proper synchronization
- File change handlers for trigger files
- Process management with background monitoring
- State transition logic

### 🎯 IMMEDIATE TASK: Claude Code Integration

**Replace simulation scripts with actual Claude Code calls:**

Currently working:

```bash
simulate-dev.sh    # Sleeps 2s, writes dev-output.md, triggers review
simulate-review.sh # Reads dev output, sleeps 2s, writes review-output.md
```

**Needs to become:**

```bash
claude-code --prompt-file templates/dev-prompt.md --context-file .claude-coordination/current-task.md
claude-code --prompt-file templates/review-prompt.md --context-file .claude-coordination/dev-session.md
```

## Development Guidelines

### Quality Standards (From Toolkit)

1. **Simplest solution that works** - Don't over-engineer
2. **Real user testing** - Actually use the system
3. **Maintain >90% test coverage** - Especially for state management
4. **Follow feature branch workflow** - Never commit directly to main

### Technical Constraints (Validated)

- **Use absolute paths**: Relative paths cause process spawning issues
- **Handle macOS sed syntax**: `sed -i ''` required for in-place editing
- **Proper process cleanup**: Prevent zombie processes
- **File locking**: Prevent coordination file corruption

## Implementation Plan

### Phase 1: Basic Claude Code Integration (This Session)

1. **Research Claude Code CLI interface** - What flags/options exist?
2. **Create prompt templates** - Based on toolkit principles
3. **Replace simulate scripts** - With actual Claude Code calls
4. **Test end-to-end** - Full workflow with real Claude Code instances

### Phase 2: Enhanced Coordination (Next Session)

1. **Add timeouts** - Prevent hanging Claude Code sessions
2. **Improve error handling** - Retry logic, better failure modes
3. **Configuration system** - Make paths and settings configurable

### Phase 3: Production Ready (Future Session)

1. **Conflict resolution** - Handle reviewer rejections
2. **Performance optimization** - Parallel processing where safe
3. **Documentation** - Setup and usage instructions

## Technical Requirements

### Must Research First

- **Claude Code CLI interface**: What's the actual command syntax?
- **Session management**: How does Claude Code handle context and output?
- **Error codes**: What exit codes does Claude Code return?

### File Structure to Create

```
templates/
├── dev-prompt.md       # Developer instance role and instructions
└── review-prompt.md    # Reviewer instance role and checklist

.claude-coordination/
├── project-context.md  # Shared project information
├── current-task.md     # Specific task being worked on
├── dev-session.md      # Developer output and decisions
└── review-feedback.md  # Reviewer analysis and suggestions
```

### Integration Points

```go
// In startDevSession() - replace this:
devScript := filepath.Join(o.workDir, "simulate-dev.sh")
o.startProcess(devScript, []string{o.workDir})

// With this (after researching actual Claude Code syntax):
claudeCmd := "claude-code" // or full path
args := []string{"--prompt-file", "templates/dev-prompt.md", "--context-file", ".claude-coordination/current-task.md"}
o.startProcess(claudeCmd, args)
```

## Success Criteria

### This Session

- [ ] Claude Code CLI research complete - understand actual interface
- [ ] Basic prompt templates created - dev and review roles defined
- [ ] First successful end-to-end run - orchestrator coordinates real Claude Code instances
- [ ] Integration test passes - state transitions work with Claude Code

### Quality Gates

- **No regression** - File watching and state management continue working
- **Real coordination** - Claude Code instances actually read each other's output
- **Error handling** - Claude Code failures properly detected and handled

## Context Preservation Notes

**DO NOT rewrite the validated orchestrator core.** Focus purely on:

1. Replacing simulation scripts with Claude Code calls
2. Creating appropriate prompt templates
3. Testing the integration

The file watching, process management, and state machine are **proven and should not be modified** unless Claude Code integration requires specific changes.

## Ready to Start

Begin by researching Claude Code's actual CLI interface to understand how to invoke it programmatically. The orchestration foundation is solid - now we need to plug in the real AI instances.
