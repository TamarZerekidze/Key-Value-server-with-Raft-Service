package kvraft

import (
	"6.5840/labgob"
	"6.5840/labrpc"
	"6.5840/raft"
	"bytes"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type Op struct {
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	ClientId       int64
	OperationsSent int64
	OperationType  string
	Key            string
	Value          string
}

type KVServer struct {
	mu      sync.Mutex
	me      int
	rf      *raft.Raft
	applyCh chan raft.ApplyMsg
	dead    int32 // set by Kill()

	maxraftstate int // snapshot if log grows this big

	// Your definitions here.
	kvMap            map[string]string
	clientsMap       map[int64]int64
	updateCh         map[int]chan Op
	lastAppliedIndex int
	persister        *raft.Persister
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	if kv.clientsMap[args.ClientId] >= args.OperationsSent {
		reply.Err = OK
		reply.Value = kv.kvMap[args.Key]
		kv.mu.Unlock()
		return
	}

	operation := Op{}
	operation.ClientId = args.ClientId
	operation.OperationsSent = args.OperationsSent
	operation.OperationType = "Get"
	operation.Key = args.Key
	operation.Value = ""
	index, _, isLeader := kv.rf.Start(operation)

	if !isLeader {
		kv.mu.Unlock()
		reply.Err = ErrWrongLeader
		return
	}

	updateCh := make(chan Op, 1)
	kv.updateCh[index] = updateCh
	kv.mu.Unlock()

	select {
	case op := <-updateCh:

		if !(op.ClientId == args.ClientId && op.OperationsSent == args.OperationsSent) {
			reply.Err = ErrWrongLeader
			return
		}
		kv.mu.Lock()
		reply.Value = kv.kvMap[args.Key] //op.Value
		reply.Err = OK
		kv.mu.Unlock()
	case <-time.After(500 * time.Millisecond):
		reply.Err = ErrWrongLeader
	}

}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	if kv.clientsMap[args.ClientId] >= args.OperationsSent {
		reply.Err = OK
		kv.mu.Unlock()
		return
	}

	operation := Op{}
	operation.ClientId = args.ClientId
	operation.OperationsSent = args.OperationsSent
	operation.OperationType = args.OperationType
	operation.Key = args.Key
	operation.Value = args.Value
	index, _, isLeader := kv.rf.Start(operation)
	if !isLeader {
		kv.mu.Unlock()
		reply.Err = ErrWrongLeader
		return
	}

	updateCh := make(chan Op, 1)
	kv.updateCh[index] = updateCh
	kv.mu.Unlock()

	select {
	case op := <-updateCh:

		if !(op.ClientId == args.ClientId && op.OperationsSent == args.OperationsSent) {
			reply.Err = ErrWrongLeader
			return
		}
		reply.Err = OK
	case <-time.After(500 * time.Millisecond):
		reply.Err = ErrWrongLeader
	}
	return
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	if kv.clientsMap[args.ClientId] >= args.OperationsSent {
		reply.Err = OK
		kv.mu.Unlock()
		return
	}

	operation := Op{}
	operation.ClientId = args.ClientId
	operation.OperationsSent = args.OperationsSent
	operation.OperationType = args.OperationType
	operation.Key = args.Key
	operation.Value = args.Value
	index, _, isLeader := kv.rf.Start(operation)

	if !isLeader {
		kv.mu.Unlock()
		reply.Err = ErrWrongLeader
		return
	}

	updateCh := make(chan Op, 1)
	kv.updateCh[index] = updateCh
	kv.mu.Unlock()

	select {
	case op := <-updateCh:

		if !(op.ClientId == args.ClientId && op.OperationsSent == args.OperationsSent) {
			reply.Err = ErrWrongLeader
			return
		}
		reply.Err = OK
	case <-time.After(500 * time.Millisecond):
		reply.Err = ErrWrongLeader
	}
	return
}

// the tester calls Kill() when a KVServer instance won't
// be needed again. for your convenience, we supply
// code to set rf.dead (without needing a lock),
// and a killed() method to test rf.dead in
// long-running loops. you can also add your own
// code to Kill(). you're not required to do anything
// about this, but it may be convenient (for example)
// to suppress debug output from a Kill()ed instance.
func (kv *KVServer) Kill() {
	atomic.StoreInt32(&kv.dead, 1)
	kv.rf.Kill()
	// Your code here, if desired.
}

func (kv *KVServer) killed() bool {
	z := atomic.LoadInt32(&kv.dead)
	return z == 1
}

func (kv *KVServer) makeSnapshot() {
	kv.mu.Lock()
	if !(kv.maxraftstate != -1 && kv.persister.RaftStateSize() >= kv.maxraftstate) {
		kv.mu.Unlock()
		return
	}
	defer kv.mu.Unlock()
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(kv.kvMap)
	e.Encode(kv.clientsMap)
	st := w.Bytes()

	kv.rf.Snapshot(kv.lastAppliedIndex, st)
	//go kv.encodeSnapshot()
}

func (kv *KVServer) resetSnapshot(snapshot []byte) {

	if len(snapshot) <= 0 {
		return
	}
	r := bytes.NewBuffer(snapshot)
	d := labgob.NewDecoder(r)

	var kvMap map[string]string
	var clientsMap map[int64]int64

	if d.Decode(&kvMap) != nil ||
		d.Decode(&clientsMap) != nil {
		DPrintf("Error")
	} else {
		kv.mu.Lock()
		defer kv.mu.Unlock()
		kv.kvMap = kvMap
		kv.clientsMap = clientsMap
	}

}

func (kv *KVServer) updatingChannel() {

	for op := range kv.applyCh {
		if op.CommandValid {
			kv.mu.Lock()
			operation := op.Command.(Op)
			clientId := operation.ClientId
			operationsSent := operation.OperationsSent
			operationType := operation.OperationType
			key := operation.Key
			value := operation.Value

			if kv.clientsMap[clientId] < operationsSent {
				if operationType == "Get" {
					operation.Value = kv.kvMap[key]
				}
				if operationType == "Append" {
					operation.Value = kv.kvMap[key] + value
					kv.kvMap[key] += value
				}
				if operationType == "Put" {
					operation.Value = value
					kv.kvMap[key] = value
				}
				kv.clientsMap[clientId] = operationsSent
			}
			kv.lastAppliedIndex = op.CommandIndex
			if kv.updateCh[op.CommandIndex] != nil {
				updateCh := kv.updateCh[op.CommandIndex]
				updateCh <- operation
				close(updateCh)
				delete(kv.updateCh, op.CommandIndex)
			}
			
			kv.mu.Unlock()

		}

		if op.SnapshotValid && !op.CommandValid {
			kv.resetSnapshot(op.Snapshot)
			kv.mu.Lock()
			kv.lastAppliedIndex = op.SnapshotIndex
			kv.mu.Unlock()
		}
		kv.makeSnapshot()

	}

}

// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots through the underlying Raft
// implementation, which should call persister.SaveStateAndSnapshot() to
// atomically save the Raft state along with the snapshot.
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
func StartKVServer(servers []*labrpc.ClientEnd, me int, persister *raft.Persister, maxraftstate int) *KVServer {
	// call labgob.Register on structures you want
	// Go's RPC library to marshall/unmarshall.
	labgob.Register(Op{})

	kv := new(KVServer)
	kv.me = me
	kv.maxraftstate = maxraftstate
	kv.persister = persister

	// You may need initialization code here.

	kv.applyCh = make(chan raft.ApplyMsg)
	kv.rf = raft.Make(servers, me, persister, kv.applyCh)

	// You may need initialization code here.
	kv.kvMap = make(map[string]string)
	kv.clientsMap = make(map[int64]int64)
	kv.updateCh = make(map[int]chan Op)
	kv.lastAppliedIndex = 0

	kv.resetSnapshot(persister.ReadSnapshot())

	go kv.updatingChannel()
	return kv
}
