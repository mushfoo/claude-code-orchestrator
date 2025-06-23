#!/bin/bash
# Run the orchestrator for Claude Code integration

echo "🚀 Starting Claude Code Orchestrator..."
echo "Working directory: $(pwd)"

# Run the orchestrator without demo mode
echo "🏃 Running orchestrator..."
go run main.go "$(pwd)"