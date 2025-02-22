# Raft: A Fault-Tolerant Replicated State Machine

This repository contains my implementation of the Raft consensus algorithm. It provides a replicated state machine that ensures consistent log replication across multiple servers.

## Raft Implementation

Raft replicates client commands in a sequential log across nodes, ensuring that all nodes execute entries in the same order. It remains fault-tolerant as long as a majority of nodes communicate. This implementation follows the extended Raft paper (excluding cluster membership changes) and is tested using the MIT lab's provided tests.

### Key/Value Storage Service
The key/value service is a replicated state machine that stores key/value pairs across multiple servers. It supports three primary RPCs:
- **Put(key, value):** Replaces the value for a key.
- **Append(key, arg):** Appends data to a key's value.
- **Get(key):** Retrieves the current value for a key.

Clients interact with the service via a Clerk, ensuring that operations are linearizable. The service is fully updated with snapshot support, allowing Raft to discard old log entries while maintaining consistency.

## Features

- **Replicated Log:** Sequential log replication across nodes.
- **Consistency:** All nodes process commands in the same order.
- **Fault Tolerance:** Recovers from node failures and network issues.
- **RPC Communication:** Nodes communicate via RPC.
- **Lab-Tested:** Validated using the MIT lab's client tests.
  
- - **Key/Value Service:**
  - Fault-tolerant, replicated state machine for key/value storage.
  - Linearizable Put, Append, and Get operations.
  - Full snapshot support for efficient log management.
