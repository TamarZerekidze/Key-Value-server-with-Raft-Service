# Raft: A Fault-Tolerant Replicated State Machine

This repository contains my implementation of the Raft consensus algorithm for the MIT 6.824 lab. It provides a replicated state machine that ensures consistent log replication across multiple servers.

## Introduction

Raft replicates client commands in a sequential log across nodes, ensuring that all nodes execute entries in the same order. It remains fault-tolerant as long as a majority of nodes communicate. This implementation follows the extended Raft paper (excluding cluster membership changes) and is tested using the MIT lab's provided tests.

## Features

- **Replicated Log:** Sequential log replication across nodes.
- **Consistency:** All nodes process commands in the same order.
- **Fault Tolerance:** Recovers from node failures and network issues.
- **RPC Communication:** Nodes communicate via RPC.
- **Lab-Tested:** Validated using the MIT lab's client tests.
