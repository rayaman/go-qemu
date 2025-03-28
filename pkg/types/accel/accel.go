package accel

type (
	VMExit        string
	Accelerator   string
	Thread        string
	KernelIrqchip string
)

const (
	KVM  Accelerator = "kvm"
	XEN  Accelerator = "xen"
	HVF  Accelerator = "hvf"
	NVMM Accelerator = "nvmm"
	// Windows
	WHPX Accelerator = "whpx"
	TCG  Accelerator = "tcg"

	Run           VMExit = "run"
	InternalError VMExit = "internal-error"
	Disable       VMExit = "disable"

	Single Thread = "single"
	Multi  Thread = "multi"

	KernelOn    KernelIrqchip = "on"
	KernelOff   KernelIrqchip = "off"
	KernelSplit KernelIrqchip = "split"
)

type Accel struct {
	Accelerator Accelerator `json:"accel,omitempty" omit:"tag"`
	// enable Xen integrated Intel graphics passthrough, default=off
	IGDPassthrough *bool `json:"igd-passthru,omitempty"`
	// controls accelerated irqchip support (default=on
	KernelIrqchip KernelIrqchip `json:"kernel-irqchip,omitempty"`
	// size of KVM shadow MMU in bytes
	KVMShadowMem uint `json:"kvm-shadow-mem,omitempty"`
	// one guest instruction per TCG translation block
	OneINSNPerTB *bool `json:"one-insn-per-tb,omitempty"`
	// enable TCG split w^x mapping
	SplitWX *bool `json:"split-wx,omitempty"`
	// TCG translation block cache size
	TBSize uint `json:"tb-size,omitempty"`
	// KVM dirty ring GFN count, default 0
	DirtyRingSize uint `json:"dirty-ring-size,omitempty"`
	// KVM Eager Page Split chunk size, default 0, disabled. ARM only
	EagerSplitSize uint `json:"eager-split-size,omitempty"`
	// enable notify VM exit and set notify window, x86 only
	NotifyVMExit VMExit `json:"notify-vm-exit,omitempty"`
	// enable notify VM exit and set notify window, x86 only
	NotifyWindow *bool `json:"notify-window,omitempty"`
	// enable multi-threaded TCG
	Thread Thread `json:"thread,omitempty"`
	// KVM device path, default /dev/kvm
	Device string `json:"device,omitempty"`
}
