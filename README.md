# Raft and Key/Value Storage Service

This repository contains my implementation of the Raft consensus algorithm along with a fully updated, fault-tolerant key/value storage service. Both components are developed as part of the MIT 6.824 distributed systems labs and are rigorously validated using the lab's provided tests.

## Overview

### Raft Consensus Algorithm
The Raft implementation provides a replicated state machine by maintaining a sequential log of client commands across multiple nodes. This ensures that every server executes the commands in the same order, preserving consistency even in the event of node failures or network issues. Following the extended Raft paper (excluding cluster membership changes), this implementation remains operational as long as a majority of nodes can communicate.

### Key/Value Storage Service
Built atop the Raft library, the key/value service stores key/value pairs in a replicated state machine. It supports three core RPC operations:
- **Put(key, value):** Replaces the current value for a key.
- **Append(key, arg):** Appends data to the existing value of a key.
- **Get(key):** Retrieves the current value for a key.

Clients use a Clerk to interact with the service, ensuring linearizable operations. With snapshot support integrated, the system can efficiently discard old log entries without sacrificing consistency.

## Features

- **Replicated Log:** Guarantees that all nodes process commands in the same order.
- **Fault Tolerance:** Remains functional even when some nodes fail or network partitions occur.
- **RPC Communication:** Facilitates efficient node-to-node interactions.
- **Snapshot Support:** Manages log size by discarding obsolete entries.
- **Lab-Tested:** Verified using the MIT lab's provided tests.

Run the MIT lab's tests to validate both the Raft library and the Key/Value service with "go test"
