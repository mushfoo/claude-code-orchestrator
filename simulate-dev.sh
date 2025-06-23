#!/bin/bash
echo "🔨 Developer Instance Starting..."
echo "Working directory: $1"

# Simulate development work (2-3 seconds)
sleep 2

# Write some "development output"
mkdir -p "$1/.claude-coordination"
echo "Development completed at $(date)" > "$1/.claude-coordination/dev-output.md"
echo "- Implemented feature X" >> "$1/.claude-coordination/dev-output.md"
echo "- Added tests" >> "$1/.claude-coordination/dev-output.md"
echo "- Ready for review" >> "$1/.claude-coordination/dev-output.md"

echo "✅ Development work complete!"

# Trigger the next stage
echo "$(date): Dev work finished" > "$1/.claude-coordination/dev-complete.trigger"
