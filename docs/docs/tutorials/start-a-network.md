---
layout: default
permalink: /tutorials/start-a-network/
---

# Start a Network

This tutorial provides a basic introduction how you can start a local network using Gosemble runtime, imported in Substrate node.

## Before you begin

Before you begin, verify that:

1. You have [installed](../development/install.md) all the repository dependencies.
2. You have [built](../development/build.md) your latest Gosemble runtime.


## Build and start the network

1. Open a terminal shell on your machine.
2. Change to the root directory to the locally cloned Gosemble repository.
3. Execute the following command:

```bash
make start-network
```

This will build the Substrate node with the Gosemble runtime wasm blob and start a network with one node.

Once the node is built, the terminal should display a similar output to this:
```bash
2023-04-20 09:00:47 Substrate Node    
2023-04-20 09:00:47 ✌️  version 4.0.0-dev-765fd435549    
2023-04-20 09:00:47 ❤️  by Substrate DevHub <https://github.com/substrate-developer-hub>, 2017-2023    
2023-04-20 09:00:47 📋 Chain specification: Development    
2023-04-20 09:00:47 🏷  Node name: real-approval-9498    
2023-04-20 09:00:47 👤 Role: AUTHORITY    
2023-04-20 09:00:47 💾 Database: RocksDb at /var/folders/4y/0ylpyqgn22g8jqpchzpm6lz80000gn/T/substrateBtT4Ur/chains/dev/db/full    
2023-04-20 09:00:47 ⛓  Native runtime: node-template-100 (node-template-1.tx1.au1)    
2023-04-20 09:00:47 🔨 Initializing Genesis block/state (state: 0x8cac…2784, header-hash: 0x3cda…df57)    
2023-04-20 09:00:47 👴 Loading GRANDPA authority set from genesis on what appears to be first startup.    
2023-04-20 09:00:47 Using default protocol ID "sup" because none is configured in the chain specs    
2023-04-20 09:00:47 🏷  Local node identity is: 12D3KooWKTKaG1R7DxRtTWGAJDAEXC91QgbvjuW2HoChuarvPVwB    
2023-04-20 09:00:47 💻 Operating system: macos    
2023-04-20 09:00:47 💻 CPU architecture: aarch64    
2023-04-20 09:00:47 📦 Highest known block at #0    
2023-04-20 09:00:47 〽️ Prometheus exporter started at 127.0.0.1:9615    
2023-04-20 09:00:47 Running JSON-RPC HTTP server: addr=127.0.0.1:9933, allowed origins=["*"]    
2023-04-20 09:00:47 Running JSON-RPC WS server: addr=127.0.0.1:9944, allowed origins=["*"]    
2023-04-20 09:00:48 🙌 Starting consensus session on top of parent 0x3cda151b8ad3c4f331710e99d76c93a6f1332fb6944274beb4942758f129df57    
2023-04-20 09:00:48 🎁 Prepared block for proposing at 1 (0 ms) [hash: 0x78f54ecfb1c9429ab0fdf79e895fe5b384996759fbad7dd080e86793cb6dd171; parent_hash: 0x3cda…df57; extrinsics (1): [0x47a9…5266]]    
2023-04-20 09:00:48 🔖 Pre-sealed block for proposal at 1. Hash now 0x73b64c2e2ebb1e36f6ce3ceae1f30db4e85ec97541cfca38f688771661283911, previously 0x78f54ecfb1c9429ab0fdf79e895fe5b384996759fbad7dd080e86793cb6dd171.    
2023-04-20 09:00:48 ✨ Imported #1 (0x73b6…3911)    
2023-04-20 09:00:50 🙌 Starting consensus session on top of parent 0x73b64c2e2ebb1e36f6ce3ceae1f30db4e85ec97541cfca38f688771661283911    
2023-04-20 09:00:50 🎁 Prepared block for proposing at 2 (0 ms) [hash: 0x8297614e7b45dde043902a55c76410ad249bdde1a34d30593a0614b0e7c8743c; parent_hash: 0x73b6…3911; extrinsics (1): [0x56e4…ec44]]    
2023-04-20 09:00:50 🔖 Pre-sealed block for proposal at 2. Hash now 0x46590bceeaf9c797c37e940b97dc7c127dfef625c540f32d3298570cdf805af1, previously 0x8297614e7b45dde043902a55c76410ad249bdde1a34d30593a0614b0e7c8743c.    
2023-04-20 09:00:50 ✨ Imported #2 (0x4659…5af1)    
2023-04-20 09:00:52 🙌 Starting consensus session on top of parent 0x46590bceeaf9c797c37e940b97dc7c127dfef625c540f32d3298570cdf805af1    
2023-04-20 09:00:52 🎁 Prepared block for proposing at 3 (1 ms) [hash: 0xc068b2a5904b34a40aeb0ee0ff64469a3879974435f36859c074542f11cacbd2; parent_hash: 0x4659…5af1; extrinsics (1): [0x2fc1…1556]]    
2023-04-20 09:00:52 🔖 Pre-sealed block for proposal at 3. Hash now 0x1f95f1d3b05ee47883cc56853029b9160f0aedf966adc874e7acc50f64a1af1f, previously 0xc068b2a5904b34a40aeb0ee0ff64469a3879974435f36859c074542f11cacbd2.    
2023-04-20 09:00:52 ✨ Imported #3 (0x1f95…af1f)    
2023-04-20 09:00:52 💤 Idle (0 peers), best: #3 (0x1f95…af1f), finalized #0 (0x3cda…df57), ⬇ 0 ⬆ 0    
2023-04-20 09:00:54 🙌 Starting consensus session on top of parent 0x1f95f1d3b05ee47883cc56853029b9160f0aedf966adc874e7acc50f64a1af1f    
2023-04-20 09:00:54 🎁 Prepared block for proposing at 4 (1 ms) [hash: 0xccb05a3ba5b0122223aceea63fdf451137f431eda74d3d5be071d033c276ad64; parent_hash: 0x1f95…af1f; extrinsics (1): [0x3933…768d]]    
2023-04-20 09:00:54 🔖 Pre-sealed block for proposal at 4. Hash now 0x499abfe622f7ba16ee2f84d93d14cfd53cfb67ad6520c2fe1d4e494feabcba08, previously 0xccb05a3ba5b0122223aceea63fdf451137f431eda74d3d5be071d033c276ad64.    
2023-04-20 09:00:54 ✨ Imported #4 (0x499a…ba08)    
2023-04-20 09:00:56 🙌 Starting consensus session on top of parent 0x499abfe622f7ba16ee2f84d93d14cfd53cfb67ad6520c2fe1d4e494feabcba08    
2023-04-20 09:00:56 🎁 Prepared block for proposing at 5 (1 ms) [hash: 0x1c360a200207e096b0b94888b35ef125636b79b7199051eb1d10e536233c1c98; parent_hash: 0x499a…ba08; extrinsics (1): [0xdfad…dc49]]    
2023-04-20 09:00:56 🔖 Pre-sealed block for proposal at 5. Hash now 0x48285138338a30e15d38ffe6d972ce295d89c32b20f393034f2aec448abf348c, previously 0x1c360a200207e096b0b94888b35ef125636b79b7199051eb1d10e536233c1c98.    
2023-04-20 09:00:56 ✨ Imported #5 (0x4828…348c)    
...
...
...
...
2023-04-20 09:01:07 💤 Idle (0 peers), best: #10 (0xa1fe…c156), finalized #7 (0x2361…27a8), ⬇ 0 ⬆ 0   
```

If the number of `finalized` blocks is increasing, this means your blockchain network is producing new blocks and successfully reaching consensus.