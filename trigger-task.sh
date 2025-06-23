#!/bin/bash
# Manually trigger a new task

echo "📝 Triggering new task..."
echo "$(date): Manual task trigger" > .claude-coordination/task-ready.trigger
echo "✅ Task triggered!"