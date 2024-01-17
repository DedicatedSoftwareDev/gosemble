---
layout: default
permalink: /development/toolchain-setup
---

Given that we use our own custom version of the TinyGo compiler, to facilitate its development process, which includes making changes, building, and executing tests, it's necessary to carry out extra steps for setting up the development environment, beyond the initial requirements specified in the [install](/development/install) page.

## Docker

There is a docker image available, which can be used to build the compiler and run tests inside a container without having to install any dependencies locally.
Run the following script to build and spin up a container:

```sh
cd tinygo
./polkawasm.sh
```

## Linux

Similar to the Dockerfile.

## MacOS
### Apple Silicon

#### Build TinyGo by using a system-wide LLVM

Depending on the TinyGo version you want to build, choose the correct version of LLVM. 
For example, TinyGo 0.25.0, requires LLVM 14. Keep in mind that even though there are some warning messages, they do not prevent builds from succeeding. Also, those warnings disappear with newer versions of LLVM and TinyGo.

```sh
brew install llvm@14
```

Make sure these environment variables are set correctly:

```sh
go env GOROOT # => /usr/local/go
go env GOPATH # => ~/go
go env GOARCH # => arm64
```

Use the Go toolchain to build TinyGo. Do not use `make`, since the `Makefile` is intended to be used with a self-built LLVM.

```sh
cd tinygo
go install
```

Make sure to include the path to the TinyGo binary in your `PATH` environment variable: 

```sh
export PATH="$GOPATH/bin:$PATH"
```

Restart the shell and verify it's working:

```sh
tinygo version
```

Install `binaryen` which is required for the wasm tests:

```sh
brew install binaryen
```

Run the tests:

```sh
# standard library packages that pass tests on darwin, linux, wasi, and windows, but take over a minute in wasi
tinygo test -target wasi compress/bzip2 crypto/dsa index/suffixarray

# standard library packages that pass tests quickly on darwin, linux, wasi, and windows
tinygo test -target wasi compress/lzw compress/zlib container/heap container/list container/ring crypto/des crypto/md5 crypto/rc4 crypto/sha1 crypto/sha256 crypto/sha512 debug/macho embed/internal/embedtest encoding encoding/ascii85 encoding/base32 encoding/base64 encoding/csv encoding/hex go/scanner hash hash/adler32 hash/crc64 hash/fnv html internal/itoa internal/profile math math/cmplx net/http/internal/ascii net/mail os path reflect sync testing testing/iotest text/scanner unicode unicode/utf16 unicode/utf8

# standard library packages that pass tests on individual platforms
tinygo test -target wasi archive/zip bytes compress/flate crypto/hmac debug/dwarf debug/plan9obj image io/ioutil mime/quotedprintable net strconv testing/fstest text/tabwriter text/template/parse

# wasi
tinygo test -target wasi ./tests/runtime_wasi

# wasm
tinygo build -size short -o wasm.wasm -target=wasm examples/wasm/export
tinygo build -size short -o wasm.wasm -target=wasm examples/wasm/main

# wasm
go test -count=1 ./tests/wasm

go test -v -count=1 ./tests/os/smoke
go test -v -count=1 ./tests/runtime
go test -v -count=1 ./tests/text/template/smoke
go test -v -count=1 ./tests/tinygotest

go test -v -count=1 -v -timeout=20m -tags "osusergo" ./builder ./cgo ./compileopts ./compiler ./interp ./transform .
```

#### Build TinyGo by using LLVM build from source

Use `make` with a self-built LLVM which has the benefit of already set up tests.

Install the necessary dependencies:

```sh
brew install cmake ninja
```

Clone and build LLVM:

```sh
make llvm-source
make llvm-build
```

Build the TinyGo compiler:

```sh
make clean
make tinygo
```

Run the tests:

```sh
make test
make smoketest
make test-corpus-wasi
make wasmtest
```
