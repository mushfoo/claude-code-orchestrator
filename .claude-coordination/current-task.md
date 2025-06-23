# Task: Add Logging Enhancements

## Description

Add structured logging enhancements to the orchestrator to improve debugging and monitoring capabilities.

## Requirements

1. Add timestamp prefixes to all log messages
2. Add log levels (INFO, WARN, ERROR)
3. Create a simple log formatter function
4. Update existing log statements to use the new format

## Example

Current log:
```
log.Printf("Starting orchestrator, watching: %s", o.coordDir)
```

Should become:
```
logInfo("Starting orchestrator, watching: %s", o.coordDir)
```

## Acceptance Criteria

- All existing log statements updated
- New log functions created (logInfo, logWarn, logError)
- Timestamps included in ISO format
- Code remains clean and maintainable