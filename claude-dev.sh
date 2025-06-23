#!/bin/bash
# Claude Code Developer Instance Wrapper

WORK_DIR="$1"
COORD_DIR="$WORK_DIR/.claude-coordination"

echo "🔨 Starting Claude Code Developer Instance..."
echo "Working directory: $WORK_DIR"

# Ensure coordination directory exists
mkdir -p "$COORD_DIR"

# Create a combined prompt that includes the role and the task
{
    cat "$WORK_DIR/templates/dev-prompt.md"
    echo ""
    echo "## Current Task"
    echo ""
    if [ -f "$COORD_DIR/current-task.md" ]; then
        cat "$COORD_DIR/current-task.md"
    else
        echo "No task found in $COORD_DIR/current-task.md"
    fi
} | claude --print --add-dir "$WORK_DIR"

# After Claude completes, trigger the next stage
echo "$(date): Dev work finished" > "$COORD_DIR/dev-complete.trigger"
echo "✅ Development phase complete!"