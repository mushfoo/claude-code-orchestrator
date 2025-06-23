# Claude Code Integration

This document describes the Claude Code integration with the orchestrator.

## Overview

The orchestrator now uses real Claude Code instances instead of simulation scripts. Two Claude Code instances are coordinated:

1. **Developer Instance** - Implements features based on tasks
2. **Reviewer Instance** - Reviews the developer's work

## File Structure

```
templates/
├── dev-prompt.md       # Developer role and instructions
└── review-prompt.md    # Reviewer role and instructions

.claude-coordination/
├── current-task.md     # Task description for developer
├── dev-session.md      # Developer's work summary
└── review-feedback.md  # Reviewer's feedback
```

## Running the Orchestrator

1. Start the orchestrator:
   ```bash
   ./run-orchestrator.sh
   ```

2. Place your task in `.claude-coordination/current-task.md`

3. Trigger the workflow:
   ```bash
   ./trigger-task.sh
   ```

## How It Works

1. **Task Ready**: When `task-ready.trigger` is modified, the orchestrator starts the developer session
2. **Developer Phase**: 
   - Runs `claude-dev.sh` which invokes Claude Code with the developer prompt
   - Claude reads the task and implements it
   - Writes summary to `dev-session.md`
   - Triggers `dev-complete.trigger`
3. **Review Phase**:
   - Runs `claude-review.sh` which invokes Claude Code with the reviewer prompt
   - Claude reads the developer's work and reviews it
   - Writes feedback to `review-feedback.md`
   - Triggers `review-complete.trigger`
4. **Cycle Complete**: Returns to idle state, ready for next task

## Testing

For testing with simulation scripts (original behavior):
```bash
go run main.go . --demo
```

## Troubleshooting

- Ensure Claude Code CLI is installed and accessible
- Check that all scripts have execute permissions
- Monitor the orchestrator logs for state transitions
- Verify trigger files are being created in `.claude-coordination/`