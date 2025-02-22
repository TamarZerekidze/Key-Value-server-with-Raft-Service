package kvraft

const (
	OK             = "OK"
	ErrNoKey       = "ErrNoKey"
	ErrWrongLeader = "ErrWrongLeader"
)

type Err string

// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	// You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	ClientId       int64
	OperationsSent int64
	OperationType  string
}

type PutAppendReply struct {
	Err Err
}

type GetArgs struct {
	Key            string
	ClientId       int64
	OperationsSent int64
}

type GetReply struct {
	Err   Err
	Value string
}
