# Gosemble

> **Warning**
> The Gosemble is in pre-production

Go implementation of Polkadot/Substrate compatible runtimes. For more details, check
the [Official Documentation](https://limechain.github.io/gosemble/)

### Quick Start

#### Prerequisites

- [Git](https://git-scm.com/downloads)
- [Go 1.19+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/install/)
- [Rust](https://docs.substrate.io/install/)

#### Clone the repository

```bash
git clone https://github.com/LimeChain/gosemble.git
cd gosemble
```

#### Pull all necessary git submodules

```bash
git submodule update --init --recursive
```

#### Build

To build a runtime, execute: 

```bash
make build
```

#### Start a local network

After the runtime is built, start a local network using Substrate host:

```bash
make start-network
```

#### Run Tests

After the Runtime is built, execute the tests with the help of [Gossamer](https://github.com/LimeChain/gossamer), which
is used to import necessary Polkadot Host functionality and interact with the Runtime.

```bash
make test_unit
make test_integration
```
