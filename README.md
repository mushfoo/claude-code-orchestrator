# Claude Code Orchestrator

A Go-based orchestration system that coordinates two Claude Code instances in a developer/reviewer workflow for automated code development and peer review.

## Overview

The Claude Code Orchestrator automates a two-phase development workflow:

1. **Developer Phase**: A Claude Code instance reads a task and implements the requested changes
2. **Review Phase**: Another Claude Code instance reviews the implementation and provides feedback

The system uses file-based triggers to coordinate between phases, ensuring a structured and traceable development process.

## Prerequisites

- Go 1.19 or higher
- Claude Code CLI installed and configured
- macOS or Linux (Windows may require adjustments)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mushfoo/claude-code-orchestrator.git
   cd claude-code-orchestrator
   ```

2. Verify Claude Code is installed:
   ```bash
   claude --version
   ```

3. Make scripts executable:
   ```bash
   chmod +x *.sh
   ```

## Directory Structure

```
claude-code-orchestrator/
├── main.go                    # Orchestrator core
├── templates/
│   ├── dev-prompt.md         # Developer instance instructions
│   └── review-prompt.md      # Reviewer instance instructions
├── .claude-coordination/      # Coordination files (auto-created)
│   ├── current-task.md       # Input: Task description
│   ├── dev-session.md        # Output: Developer's work summary
│   ├── review-feedback.md    # Output: Reviewer's feedback
│   └── *.trigger             # Trigger files for state transitions
└── Scripts:
    ├── run-orchestrator.sh   # Start the orchestrator
    ├── trigger-task.sh       # Manually trigger a task
    ├── claude-dev.sh         # Developer wrapper (used internally)
    └── claude-review.sh      # Reviewer wrapper (used internally)
```

## Quick Start

### Step 1: Start the Orchestrator

```bash
./run-orchestrator.sh
```

You should see:
```
🚀 Starting orchestrator, watching: .claude-coordination
📝 To start a task: echo 'start' > .claude-coordination/task-ready.trigger
👂 Watching for file changes... (press Ctrl+C to exit)
```

### Step 2: Create a Task

Create or edit `.claude-coordination/current-task.md` with your task:

```markdown
# Task: Add User Authentication

## Description
Implement a basic user authentication system with login and logout functionality.

## Requirements
1. Create a User model with email and password
2. Add login endpoint that returns a session token
3. Add logout endpoint that invalidates the session
4. Include appropriate error handling
5. Write unit tests for the authentication logic
```

### Step 3: Trigger the Workflow

```bash
./trigger-task.sh
```

### Step 4: Monitor Progress

The orchestrator will show the workflow progress:
```
📁 File changed: task-ready.trigger
🔨 Starting Developer session...
[Claude Code developer output]
✅ Development phase complete!

📁 File changed: dev-complete.trigger  
👀 Starting Review session...
[Claude Code reviewer output]
✅ Review phase complete!

✅ Cycle complete, returning to idle
```

### Step 5: Review Results

Check the output files:
- `.claude-coordination/dev-session.md` - Developer's implementation summary
- `.claude-coordination/review-feedback.md` - Reviewer's feedback

## Detailed Workflow

### 1. Idle State
- Orchestrator watches for changes to `task-ready.trigger`
- You prepare task in `current-task.md`

### 2. Developer Phase (StateDevRunning)
- Triggered by: `task-ready.trigger` modification
- Executes: `claude-dev.sh`
- Claude Code reads the task and implements it
- Outputs summary to `dev-session.md`
- Triggers: `dev-complete.trigger`

### 3. Review Phase (StateReviewRunning)
- Triggered by: `dev-complete.trigger` modification
- Executes: `claude-review.sh`
- Claude Code reviews the developer's work
- Outputs feedback to `review-feedback.md`
- Triggers: `review-complete.trigger`

### 4. Return to Idle
- Triggered by: `review-complete.trigger` modification
- Ready for next task

## File-Based Triggers

The system uses trigger files that you can manually control:

```bash
# Start a new task (if orchestrator is in Idle state)
echo "start" > .claude-coordination/task-ready.trigger

# Skip to review (if orchestrator is in DevRunning state)
echo "done" > .claude-coordination/dev-complete.trigger

# Complete review (if orchestrator is in ReviewRunning state)
echo "done" > .claude-coordination/review-complete.trigger
```

## Advanced Usage

### Running with Demo Mode

To see the orchestrator work with simulated scripts:
```bash
go run main.go . --demo
```

### Custom Task Templates

Edit the prompt templates in `templates/` to customize the behavior:
- `dev-prompt.md` - Modify developer instructions
- `review-prompt.md` - Modify review criteria

### Manual Orchestrator Control

Run the orchestrator directly:
```bash
go run main.go /path/to/project
```

## Troubleshooting

### Common Issues

1. **"Script not found" errors**
   - Ensure all `.sh` files are executable: `chmod +x *.sh`
   - Run from the project root directory

2. **Claude Code not responding**
   - Verify Claude Code is installed: `which claude`
   - Check Claude Code works standalone: `echo "test" | claude --print`

3. **State machine stuck**
   - Check current state in logs
   - Manually trigger next phase using trigger files
   - Restart orchestrator if needed

### Debug Mode

View detailed logs by checking the orchestrator output. Each state transition and file change is logged with emojis for easy tracking:
- 🚀 Starting
- 📁 File changes
- 🔨 Developer phase
- 👀 Review phase
- ✅ Success
- ❌ Errors

## Development

### Running Tests

```bash
# Run with test script (rapid triggers)
./test-rapid.sh
```

### Architecture

The orchestrator uses:
- `fsnotify` for file watching
- Mutex-protected state machine
- Process management with proper cleanup
- Absolute paths for reliability

### State Machine

```
Idle → DevRunning → ReviewRunning → Idle
         ↓              ↓
       Error ←────────Error
```

## Contributing

1. Keep the validated core orchestrator logic intact
2. Test changes with both real Claude Code and simulation scripts
3. Update documentation for any new features
4. Follow the existing code style and patterns

## License

[Your license here]

## Support

For issues or questions:
- Check `CLAUDE_INTEGRATION.md` for technical details
- Review orchestrator logs for state information
- Open an issue with full error output and steps to reproduce