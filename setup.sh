#!/bin/bash

# Setup script for testing the file watcher spike

echo "🔧 Setting up file watcher spike test environment..."

# Create test directory structure
mkdir -p spike-test
cd spike-test

# Create simulation scripts that mimic Claude Code behavior
cat > simulate-dev.sh << 'EOF'
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
EOF

cat > simulate-review.sh << 'EOF'
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
EOF

# Make scripts executable
chmod +x simulate-dev.sh
chmod +x simulate-review.sh

# Create go.mod for the spike
cat > go.mod << 'EOF'
module claude-coordinator-spike

go 1.21

require github.com/fsnotify/fsnotify v1.7.0

require golang.org/x/sys v0.4.0 // indirect
EOF

# Copy the main.go from the artifact (user will need to do this)
echo "📋 Next steps:"
echo "1. Copy the main.go code from the artifact into this directory"
echo "2. Run: go mod tidy"
echo "3. Test the spike: go run main.go ."
echo ""
echo "🧪 Test scenarios to try:"
echo "- Normal workflow: task → dev → review → complete"
echo "- Process interruption: kill simulate-dev.sh mid-run"
echo "- Rapid triggers: trigger multiple events quickly"
echo "- Error handling: make simulate-dev.sh exit with error"
echo ""
echo "📁 Directory structure created:"
find . -type f -exec ls -la {} \;
