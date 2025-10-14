---
name: Bug Report
about: Create a report to help us improve htmlemail
title: '[BUG] '
labels: ['bug', 'needs-triage']
assignees: ''
---

## Bug Description

A clear and concise description of what the bug is.

## Steps to Reproduce

1. Go to '...'
2. Create template with '...'
3. Call method '...'
4. See error

## Expected Behavior

A clear and concise description of what you expected to happen.

## Actual Behavior

A clear and concise description of what actually happened.

## Code Sample

Please provide a minimal code sample that reproduces the issue:

```go
package main

import "github.com/gusdeyw/htmlemail"

func main() {
    // Your code that reproduces the bug
    template := htmlemail.NewEmailTemplateFromString("...")
    // ...
}
```

## Error Message

If applicable, paste the full error message:

```
Error message here
```

## Environment

- **Go version**: [e.g. go1.21.0]
- **OS**: [e.g. Windows 11, macOS 14.0, Ubuntu 22.04]
- **htmlemail version**: [e.g. v1.0.0]
- **Architecture**: [e.g. amd64, arm64]

## Template Content

If the issue involves a specific template, please provide it:

```html
<!-- Your HTML template here -->
```

## Additional Context

Add any other context about the problem here, such as:
- Does this happen consistently or intermittently?
- Did this work in a previous version?
- Are there any workarounds you've found?

## Possible Solution

If you have an idea for how to fix this bug, please describe it here.

## Related Issues

Link any related issues: #123