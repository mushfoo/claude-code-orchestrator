#!/bin/bash
echo "🚀 Testing rapid fire triggers..."

# Fire multiple triggers quickly
echo "rapid 1 $(date)" > .claude-coordination/task-ready.trigger
sleep 0.2
echo "rapid 2 $(date)" > .claude-coordination/dev-complete.trigger  
sleep 0.2
echo "rapid 3 $(date)" > .claude-coordination/review-complete.trigger

echo "✅ Rapid triggers completed"
