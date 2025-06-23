# Claude Code Orchestrator

A Go-based orchestration system that coordinates two Claude Code instances in a developer/reviewer workflow for automated code development and peer review.

## Overview

The Claude Code Orchestrator automates a structured two-phase development workflow:

1. **Developer Phase**: A Claude Code instance reads a task description and implements the requested changes
2. **Review Phase**: Another Claude Code instance reviews the implementation and provides constructive feedback

The system uses a file-based trigger mechanism and state machine to coordinate between phases, ensuring a traceable and organized development process.

## Prerequisites

- **Go 1.21 or higher** (as specified in go.mod)
- **Claude Code CLI** installed and configured with valid API access
- **Unix-like environment** (macOS/Linux - Windows may require script modifications)
- **Git** for version control integration

## Installation

1. **Clone and navigate to the repository:**
   ```bash
   git clone https://github.com/mushfoo/claude-code-orchestrator.git
   cd claude-code-orchestrator
   ```

2. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

3. **Verify Claude Code installation:**
   ```bash
   claude --version
   ```

4. **Make all scripts executable:**
   ```bash
   chmod +x *.sh
   ```

5. **Test the installation:**
   ```bash
   go run main.go . --demo
   ```

## Project Structure

```
claude-code-orchestrator/
├── main.go                    # Core orchestrator with state machine
├── go.mod                     # Go module definition  
├── go.sum                     # Go dependency checksums
├── LICENSE                    # Project license
├── README.md                  # This documentation
├── templates/                 # Prompt templates for Claude instances
│   ├── dev-prompt.md         # Developer instance instructions
│   └── review-prompt.md      # Reviewer instance instructions
├── .claude-coordination/      # Coordination directory (auto-created)
│   ├── current-task.md       # Input: Task description
│   ├── dev-session.md        # Output: Developer's work summary
│   ├── review-feedback.md    # Output: Reviewer's feedback
│   ├── task-ready.trigger    # Trigger: Start development phase
│   ├── dev-complete.trigger  # Trigger: Start review phase
│   └── review-complete.trigger # Trigger: Complete cycle
├── claude-dev.sh             # Developer instance wrapper script
├── claude-review.sh          # Reviewer instance wrapper script
├── run-orchestrator.sh       # Orchestrator startup script
├── trigger-task.sh           # Manual task trigger utility
└── test-rapid.sh             # Rapid testing script for development
```

## Quick Start

### Step 1: Start the Orchestrator

```bash
./run-orchestrator.sh
```

Expected output:
```
🚀 Starting Claude Code Orchestrator...
Working directory: /path/to/your/project
🏃 Running orchestrator...
🚀 Starting orchestrator, watching: .claude-coordination
📝 To start a task: echo 'start' > .claude-coordination/task-ready.trigger
👂 Watching for file changes... (press Ctrl+C to exit)
```

### Step 2: Create a Task Description

Create or edit `.claude-coordination/current-task.md` with your development task:

```markdown
# Task: Add User Authentication

## Description
Implement a basic user authentication system with login and logout functionality.

## Requirements
1. Create a User model with email and password fields
2. Add login endpoint that validates credentials and returns a session token
3. Add logout endpoint that invalidates the session token
4. Include appropriate error handling and validation
5. Write comprehensive unit tests for the authentication logic
6. Update documentation to reflect the new authentication flow

## Acceptance Criteria
- All tests must pass
- Code follows existing project conventions
- Security best practices are implemented
- API endpoints are properly documented
```

### Step 3: Trigger the Development Workflow

```bash
./trigger-task.sh
```

Alternative manual trigger:
```bash
echo "$(date): Starting new task" > .claude-coordination/task-ready.trigger
```

### Step 4: Monitor the Automated Workflow

The orchestrator will progress through states automatically:

```
📁 File changed: task-ready.trigger (current state: Idle)
🔨 Starting Developer session...
🔍 Dev script path: /path/to/project/claude-dev.sh
[Claude Code developer instance output]
✅ Development phase complete!

📁 File changed: dev-complete.trigger (current state: DevRunning)  
👀 Starting Review session...
🔍 Review script path: /path/to/project/claude-review.sh
[Claude Code reviewer instance output]
✅ Review phase complete!

📁 File changed: review-complete.trigger (current state: ReviewRunning)
✅ Cycle complete, returning to idle
```

### Step 5: Review Results and Outputs

Check the generated files in `.claude-coordination/`:
- **`dev-session.md`** - Developer's implementation summary and decisions
- **`review-feedback.md`** - Reviewer's code review and recommendations
- **Trigger files** - Track the workflow state transitions

## Detailed Workflow

### State Machine Overview

The orchestrator implements a simple state machine with the following transitions:

```
┌─────────────┐     task-ready.trigger     ┌──────────────────┐
│    Idle     │ ───────────────────────────▶│   DevRunning     │
│             │                            │                  │
└─────────────┘                            └──────────────────┘
       ▲                                              │
       │                                              │ dev-complete.trigger
       │ review-complete.trigger                      ▼
       │                                   ┌──────────────────┐
       └───────────────────────────────────│  ReviewRunning   │
                                          │                  │
                                          └──────────────────┘
```

### 1. Idle State
- **Orchestrator watches** for changes to `task-ready.trigger`
- **User prepares** task description in `current-task.md`
- **State**: Ready to begin new development cycle

### 2. Developer Phase (StateDevRunning)
- **Triggered by**: `task-ready.trigger` file modification
- **Executes**: `claude-dev.sh` script with project directory
- **Process**: 
  - Combines `dev-prompt.md` template with `current-task.md`
  - Launches Claude Code with developer role and task context
  - Claude Code analyzes, implements, and tests the requested changes
- **Outputs**: Implementation summary to `dev-session.md`
- **Completion**: Automatically triggers `dev-complete.trigger`

### 3. Review Phase (StateReviewRunning)
- **Triggered by**: `dev-complete.trigger` file modification
- **Executes**: `claude-review.sh` script with project directory
- **Process**:
  - Combines `review-prompt.md` template with `dev-session.md`
  - Launches Claude Code with reviewer role and development context
  - Claude Code reviews changes, tests, and code quality
- **Outputs**: Review feedback to `review-feedback.md`
- **Completion**: Automatically triggers `review-complete.trigger`

### 4. Return to Idle
- **Triggered by**: `review-complete.trigger` file modification
- **Process**: State machine returns to Idle
- **Ready**: Available for next development cycle

## File-Based Trigger System

The system uses trigger files for state coordination. You can manually control the workflow:

```bash
# Start a new development task (when in Idle state)
echo "$(date): New task started" > .claude-coordination/task-ready.trigger

# Force skip to review phase (when in DevRunning state)
echo "$(date): Development complete" > .claude-coordination/dev-complete.trigger

# Complete the review cycle (when in ReviewRunning state)
echo "$(date): Review complete" > .claude-coordination/review-complete.trigger
```

### Trigger File Behavior
- **File watching**: Uses `fsnotify` to detect file modifications
- **Content agnostic**: Any write operation to trigger files activates the transition
- **Automatic creation**: Trigger files are created automatically if missing
- **Manual override**: Allows debugging and manual workflow control

## Advanced Usage

### Demo Mode for Testing

Test the orchestrator with simulated workflow (no actual Claude Code calls):
```bash
go run main.go . --demo
```

Demo mode automatically:
1. Triggers a new task after 1 second
2. Simulates development completion after 3 seconds  
3. Simulates review completion after 3 more seconds
4. Returns to idle state

### Custom Prompt Templates

Customize Claude Code behavior by editing templates in `templates/`:

**`templates/dev-prompt.md`** - Developer instance instructions:
- Modify development approach and coding standards
- Adjust output format for `dev-session.md`
- Add project-specific context and constraints

**`templates/review-prompt.md`** - Reviewer instance instructions:
- Customize review criteria and checklists
- Modify feedback format for `review-feedback.md`
- Add security, performance, or style requirements

### Direct Orchestrator Execution

Run the orchestrator with custom working directory:
```bash
go run main.go /path/to/your/project
```

### Rapid Testing and Development

Use the rapid testing script during development:
```bash
./test-rapid.sh
```

This script fires all three triggers in quick succession (0.2s intervals) to test state transitions.

### Integration with CI/CD

The orchestrator can be integrated into automated workflows:

```yaml
# Example GitHub Actions integration
- name: Run Orchestrated Development
  run: |
    ./run-orchestrator.sh &
    ORCH_PID=$!
    echo "${{ github.event.issue.body }}" > .claude-coordination/current-task.md
    ./trigger-task.sh
    # Wait for completion or timeout
    sleep 300
    kill $ORCH_PID
```

## Troubleshooting

### Common Issues

#### 1. Script Execution Problems
**Error**: `"Script not found"` or permission denied errors
- **Solution**: Make all scripts executable: `chmod +x *.sh`
- **Check**: Run from the project root directory
- **Verify**: Scripts exist and have correct paths

#### 2. Claude Code Integration Issues  
**Error**: Claude Code not responding or command not found
- **Check installation**: `which claude` and `claude --version`
- **Test standalone**: `echo "Hello" | claude --print`
- **Verify authentication**: Claude Code API access configured
- **Check PATH**: Claude Code CLI in system PATH

#### 3. State Machine Problems
**Error**: Orchestrator stuck in one state
- **Check logs**: Review orchestrator output for error messages
- **Manual trigger**: Force state transition with trigger files
- **Process check**: Kill any hung Claude Code processes
- **Restart**: Stop and restart the orchestrator

#### 4. File System Issues
**Error**: Coordination directory or trigger files problems
- **Permissions**: Ensure write access to project directory
- **Directory**: Verify `.claude-coordination/` exists and is writable
- **Space**: Check available disk space
- **Cleanup**: Remove corrupted trigger files and restart

#### 5. Go Module and Dependency Issues
**Error**: Build or runtime errors
```bash
# Clean and rebuild dependencies
go mod tidy
go mod download
go clean -cache

# Test compilation
go build -o orchestrator main.go
```

### Debug Mode and Logging

The orchestrator provides detailed logging with emoji indicators:

- 🚀 **Starting** - Orchestrator initialization
- 📁 **File changes** - Trigger file modifications detected  
- 🔨 **Developer phase** - Development instance running
- 👀 **Review phase** - Review instance running
- ✅ **Success** - Phase completion or successful operations
- ❌ **Errors** - Error conditions and failures
- 🔍 **Debug info** - Script paths and detailed process information

### Diagnostic Commands

```bash
# Check orchestrator state and processes
ps aux | grep "go run main.go"
ps aux | grep claude

# Monitor trigger files
watch -n 1 'ls -la .claude-coordination/*.trigger'

# View recent coordination files
ls -la .claude-coordination/
tail .claude-coordination/*.md

# Test file watching
echo "test $(date)" > .claude-coordination/task-ready.trigger
```

### Recovery Procedures

**Complete Reset**:
```bash
# Stop orchestrator (Ctrl+C)
# Clean coordination directory
rm -rf .claude-coordination/
# Restart orchestrator
./run-orchestrator.sh
```

**State Reset Only**:
```bash
# Keep files, reset triggers only
rm .claude-coordination/*.trigger
# Orchestrator will recreate trigger files
```

## Development

### Testing the Orchestrator

```bash
# Run rapid state transition testing
./test-rapid.sh

# Run demo mode for functional testing
go run main.go . --demo

# Manual state testing
./run-orchestrator.sh
# In another terminal:
./trigger-task.sh
```

### Building and Installation

```bash
# Build binary
go build -o orchestrator main.go

# Install dependencies
go mod tidy

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o orchestrator-linux main.go
GOOS=darwin GOARCH=amd64 go build -o orchestrator-macos main.go
```

### Architecture Details

#### Core Components

**State Machine (`ProcessState`)**:
- `StateIdle` - Waiting for tasks
- `StateDevRunning` - Developer instance active  
- `StateReviewRunning` - Reviewer instance active
- `StateError` - Error condition requiring intervention

**Orchestrator Structure**:
- **File Watcher** - `fsnotify.Watcher` monitors coordination directory
- **Process Management** - Handles Claude Code instance lifecycle
- **Thread Safety** - Mutex-protected state transitions
- **Error Handling** - Robust error recovery and logging

**File System Integration**:
- **Absolute Paths** - Ensures reliability across environments
- **Atomic Operations** - File-based triggers prevent race conditions
- **Auto-creation** - Missing directories and files created automatically

#### Dependencies

```go
require (
    github.com/fsnotify/fsnotify v1.7.0 // File system event notifications
    golang.org/x/sys v0.4.0             // System-specific functionality
)
```

### Extended State Machine

```
                    ┌─────────────┐
                    │    Error    │◄──────────────────────┐
                    │   (Manual   │                       │
                    │  Recovery)  │                       │
                    └─────────────┘                       │
                                                          │
    ┌─────────────┐     task-ready.trigger     ┌──────────────────┐
    │    Idle     │ ───────────────────────────▶│   DevRunning     │
    │             │                            │  (claude-dev.sh) │──┐
    └─────────────┘                            └──────────────────┘  │
           ▲                                              │           │
           │                                              │           │
           │ review-complete.trigger                      ▼           │
           │                                   ┌──────────────────┐  │ Error
           └───────────────────────────────────│  ReviewRunning   │  │ Conditions
                                              │(claude-review.sh)│◄─┘
                                              └──────────────────┘
```

## Contributing

### Development Guidelines

1. **Preserve Core Logic**: The validated state machine and file watching mechanisms are battle-tested
2. **Test Thoroughly**: Run both demo mode and real Claude Code integration tests
3. **Update Documentation**: Keep README and code comments current
4. **Follow Patterns**: Maintain existing code style and logging format
5. **Error Handling**: Ensure robust error recovery and user feedback

### Code Standards

- **Logging**: Use emoji prefixes for log levels (🚀 🔨 👀 ✅ ❌)
- **File Operations**: Always use absolute paths and check file existence
- **Concurrency**: Protect shared state with appropriate mutexes
- **Process Management**: Clean up child processes and handles properly

### Testing Contributions

```bash
# Test your changes
go run main.go . --demo              # Functional test
./test-rapid.sh                      # State transition test
go build && ./orchestrator .         # Build verification

# Test with real Claude Code
./run-orchestrator.sh
echo "Test task" > .claude-coordination/current-task.md
./trigger-task.sh
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support and Community

### Getting Help

1. **Check Documentation**: Review this README and inline code comments
2. **Enable Debug Logging**: Monitor orchestrator output for detailed state information  
3. **Test in Isolation**: Use demo mode to verify orchestrator behavior
4. **Check Dependencies**: Ensure Go, Claude Code CLI, and file permissions

### Reporting Issues

When reporting problems, include:
- Full error output and logs
- Steps to reproduce the issue
- System information (OS, Go version, Claude Code version)
- Contents of `.claude-coordination/` directory

### Feature Requests

The orchestrator is designed to be extensible. Consider:
- Additional trigger mechanisms
- Enhanced logging and monitoring
- Integration with other AI coding tools
- Custom workflow templates