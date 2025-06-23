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

## License

MIT License - Copyright (c) 2025 mushfoo

See the [LICENSE](LICENSE) file for full details.