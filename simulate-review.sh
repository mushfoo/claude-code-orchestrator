#!/bin/bash
echo "👀 Review Instance Starting..."
echo "Working directory: $1"

# Read the dev output
if [ -f "$1/.claude-coordination/dev-output.md" ]; then
    echo "📖 Reading development output..."
    cat "$1/.claude-coordination/dev-output.md"
fi

# Simulate review work (2-3 seconds)
sleep 2

# Write review feedback
echo "Review completed at $(date)" > "$1/.claude-coordination/review-output.md"
echo "- Code looks good" >> "$1/.claude-coordination/review-output.md"
echo "- Tests are adequate" >> "$1/.claude-coordination/review-output.md"
echo "- Approved for merge" >> "$1/.claude-coordination/review-output.md"

echo "✅ Review complete!"

# Trigger cycle completion
echo "$(date): Review finished" > "$1/.claude-coordination/review-complete.trigger"
