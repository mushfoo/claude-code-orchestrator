#!/bin/bash
# Claude Code Reviewer Instance Wrapper

WORK_DIR="$1"
COORD_DIR="$WORK_DIR/.claude-coordination"

echo "👀 Starting Claude Code Reviewer Instance..."
echo "Working directory: $WORK_DIR"

# Create a combined prompt that includes the role and the dev work
{
    cat "$WORK_DIR/templates/review-prompt.md"
    echo ""
    echo "## Developer's Work Summary"
    echo ""
    if [ -f "$COORD_DIR/dev-session.md" ]; then
        cat "$COORD_DIR/dev-session.md"
    else
        echo "No development session found in $COORD_DIR/dev-session.md"
    fi
} | claude --print --add-dir "$WORK_DIR"

# After Claude completes, trigger cycle completion
echo "$(date): Review finished" > "$COORD_DIR/review-complete.trigger"
echo "✅ Review phase complete!"