---
layout: default
title: Build
permalink: /development/build
---

By utilizing our [Toolchain](../../overview/toolchain), there are currently two options to choose from for the GC
implementation. Modify the `GC` environment variable to switch between them.

### Extalloc GC

It works with the host's external allocator as per specification.

```bash
make build
```

### Conservative GC

It is used only for **development** and **testing** and works by using a different heap base offset from the allocator's
one (as a workaround), so the GC can use a separate heap region for its allocations and not interfere with the
allocator's region.

```bash
GC="conservative" make build
```