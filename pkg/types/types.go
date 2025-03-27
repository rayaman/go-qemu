package types

import "fmt"

type (
	Preallocation   string
	Compat          string
	CompressionType string
	Format          string
	System          string
	Accelerator     string
	VMExit          string
	Thread          string
	KernelIrqchip   string
	Flag            string
)

// converts bool to on/off format
var SW = map[bool]string{
	true:  "on",
	false: "off",
}
var (
	On *bool = func(b bool) *bool {
		return &b
	}(true)
	Off *bool = func(b bool) *bool {
		return &b
	}(false)
)

type Param interface {
	Expand() string
}

var register = map[string]any{}

var (
	// No pre-allocation
	OFF Preallocation = "off"
	// Allocates qcow2 metadata, and it's still a sparse image.
	METADATA Preallocation = "metadata"
	// Uses posix_fallocate() to "allocate blocks and marking them as uninitialized", and is relatively faster than writing out zeroes to a file:
	FALLOC Preallocation = "falloc"
	// Allocates zeroes and makes a non-sparse image.
	FULL Preallocation = "full"

	Compat_0_10 Compat = "0.10"
	Compat_1_1  Compat = "1.1"

	Zlib CompressionType = "zlib"
	Zstd CompressionType = "zstd"

	AARCH64       System = "aarch64"
	AARCH64W      System = "aarch64w"
	ALPHA         System = "alpha"
	ALPHAW        System = "alphaw"
	ARM           System = "arm"
	ARMW          System = "armw"
	AVR           System = "avr"
	AVRW          System = "avrw"
	HPPA          System = "hppa"
	HPPAW         System = "hppaw"
	I386          System = "i386"
	I386W         System = "i386w"
	LOONGARCH64   System = "loongarch64"
	LOONGARCH64W  System = "loongarch64w"
	M68K          System = "m68k"
	M68KW         System = "m68kw"
	MICROBLAZE    System = "microblaze"
	MICROBLAZEEL  System = "microblazeel"
	MICROBLAZEELW System = "microblazeelw"
	MICROBLAZEW   System = "microblazew"
	MIPS          System = "mips"
	MIPS64        System = "mips64"
	MIPS64EL      System = "mips64el"
	MIPS64ELW     System = "mips64elw"
	MIPS64W       System = "mips64w"
	MIPSEL        System = "mipsel"
	MIPSELW       System = "mipselw"
	MIPSW         System = "mipsw"
	OR1K          System = "or1k"
	OR1KW         System = "or1kw"
	PPC           System = "ppc"
	PPC64         System = "ppc64"
	PPC64W        System = "ppc64w"
	PPCW          System = "ppcw"
	RISCV32       System = "riscv32"
	RISCV32W      System = "riscv32w"
	RISCV64       System = "riscv64"
	RISCV64W      System = "riscv64w"
	RX            System = "rx"
	RXW           System = "rxw"
	S390X         System = "s390x"
	S390XW        System = "s390xw"
	SH4           System = "sh4"
	SH4EB         System = "sh4eb"
	SH4EBW        System = "sh4ebw"
	SH4W          System = "sh4w"
	SPARC         System = "sparc"
	SPARC64       System = "sparc64"
	SPARC64W      System = "sparc64w"
	SPARCW        System = "sparcw"
	TRICORE       System = "tricore"
	TRICOREW      System = "tricorew"
	X86_64        System = "x86_64"
	X86_64W       System = "x86_64w"
	XTENSA        System = "xtensa"
	XTENSAEB      System = "xtensaeb"
	XTENSAEBW     System = "xtensaebw"
	XTENSAW       System = "xtensaw"

	Set Flag = "flag-on"
)

func GetTypes() map[string]any {
	return register
}

func RegisterType(t string, i any) error {
	if _, ok := register[t]; !ok {
		register[t] = i
		return nil
	}
	return fmt.Errorf("type already registered")
}
