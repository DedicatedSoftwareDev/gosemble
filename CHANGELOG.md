
0.0.0 (TODO)
---
- **toolchain**
  - add separate Dockerfile and build script for building TinyGo with prebuild LLVM and dependencies for faster
    builds.
  - add new target similar to Rust's `wasm32-unknown-unknown` aimed to support Polkadot's Wasm MVP.
  - linker flag to declare memory as imported.
  - linker flags to export `__heap_base`, `__data_end` globals.
  - linker flags to export `__indirect_function_table`.
  - disable the scheduler to remove the support of *goroutines* and channels (and JS/WASI exports).
  - remove the unsupported features by Wasm MVP (bulk memory operations, lang. ext) and add implementation
    of `memmove`, `memset`, `memcpy`, use opt flag as part of the target.
  - remove the exported allocation functions.
  - implement GC that can work with external memory allocator.
- **scale codec**
  - implement SCALE codec with minimal reflection
- **runtime apis**  
  - Core API.
    - `Core_version`
    - `Core_execute_block`
    - `Core_initialize_block`
  - Metadata API.
    - `Metadata_metadata`
  - BlockBuilder API.
    - `BlockBuilder_apply_extrinsic`
    - `BlockBuilder_finalize_block`
    - `BlockBuilder_inherent_extrinsics`
    - `BlockBuilder_check_inherents`
  - TaggedTransactionQueue API.
    - `TaggedTransactionQueue_validate_transaction`
  - OffchainWorker API.
    - `OffchainWorkerApi_offchain_worker`
  - Grandpa API.
    - `GrandpaApi_grandpa_authorities`
  - SessionKeys API.
    - `SessionKeys_generate_session_keys`
    - `SessionKeys_decode_session_keys`
  - AccountNonce API.
    - `AccountNonceApi_account_nonce`
  - TransactionPayment API.
    - `TransactionPaymentApi_query_info`
    - `TransactionPaymentApi_query_fee_details`
    - `TransactionPaymentCallApi_query_call_info`
    - `TransactionPaymentCallApi_query_call_fee_details`
  - Aura API.
    - `AuraApi_slot_duration`
    - `AuraApi_authorities`
- **development & tests**
  - setup development and test environment by utilizing Gossamer and Substrate hosts.



