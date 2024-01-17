---
layout: default
title: File structure
permalink: /development/file-structure
---

* `build` - the output directory for the compiled Wasm file.
* `config` - configuration of the used runtime modules (pallets).
* `constants` - constants used in the runtime.
* `env` - stubs for the host-provided functions.
* `execution` - runtime execution logic.
* `frame` - runtime modules (pallets).
* `primitives` - runtime primitives.
* `runtime` - runtime entry point and tests.
* `utils` - utility functions.
* `tinygo` - submodule for the TinyGo compiler, used for WASM compilation.
* `goscale` - submodule for the SCALE codec.
* `gossamer` - submodule for the Gossamer host, used during development and for running tests.
* `substrate` - submodule for the Substrate host, used for running a network.