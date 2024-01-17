---
layout: default
permalink: /development/test
---

Currently, the project contains unit and integration tests. Integration tests use [Gossamer](https://github.com/LimeChain/gossamer), which
imports all the necessary Host functions and interacts with the Runtime.

```bash
make test
```

or

```bash
make test_unit
make test_integration
```

### Debug

To aid the debugging process, there is a set of imported functions that can be called within the Runtime to log messages.

```go
func Critical(message string) // logs and aborts the execution
func Warn(message string)
func Info(message string)
func Debug(message string)
func Trace(message string)
```