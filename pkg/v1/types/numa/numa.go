package numa

type (
	Hierarchy     string
	Latency       string
	Associativity string
	Policy        string
)

const (
	H_Mem    Hierarchy = "memory"
	H_First  Hierarchy = "first-level"
	H_Second Hierarchy = "second-level"
	H_Third  Hierarchy = "third-level"

	L_Access Latency = "access-latency"
	L_Read   Latency = "read-latency"
	L_Write  Latency = "write-latency"

	A_None    Associativity = "none"
	A_Direct  Associativity = "direct"
	A_Complex Associativity = "complex"

	P_None      Policy = "none"
	P_WriteBack Policy = "write-back"
	P_WriteThru Policy = "write-thrugh"
)

type CPUS struct {
	First uint `json:"first,omitempty"`
	Last  uint `json:"last,omitempty"`
}
type Node struct {
	Memory    uint `json:"mem,omitempty" mexclusive:"memdev"`
	MemoryDev uint `json:"memdev,omitempty" mexclusive:"mem"`
	Cpus      CPUS `json:"cpus,omitempty"`
	NodeID    uint `json:"nodeid,omitempty"`
	Initiator uint `json:"initiator,omitempty"`
}
type CPU struct {
	NodeID   uint `json:"node-id,omitempty"`
	SocketID uint `json:"socket-id,omitempty"`
	CoreID   uint `json:"core-id,omitempty"`
	ThreadID uint `json:"thread-id,omitempty"`
}

type HMATLB struct {
	Initiator uint      `json:"initiator,omitempty"`
	Target    uint      `json:"target,omitempty"`
	Hierarchy Hierarchy `json:"hierarchy,omitempty"`
	DataType  Latency   `json:"data-type,omitempty"`
	Latency   uint      `json:"latency,omitempty"`
	Bandwidth uint      `json:"bandwidth,omitempty"`
}
type HMATCache struct {
	NodeID        uint          `json:"node-id,omitempty"`
	Size          uint          `json:"size,omitempty"`
	Level         uint          `json:"level,omitempty"`
	Associativity Associativity `json:"associativity,omitempty"`
	Policy        Policy        `json:"policy,omitempty"`
	Line          uint          `json:"line,omitempty"`
}
type Dist struct {
	Src uint `json:"src,omitempty"`
	Dst uint `json:"dst,omitempty"`
	Val uint `json:"val,omitempty"`
}
type Numa struct {
	Nodes     []Node      `json:"node,omitempty"`
	Dist      []Dist      `json:"dist,omitempty"`
	Cpu       []CPU       `json:"cpu,omitempty"`
	HMATLB    []HMATLB    `json:"hmat-lb,omitempty"`
	HMATCache []HMATCache `json:"hmat-cache,omitempty"`
}
