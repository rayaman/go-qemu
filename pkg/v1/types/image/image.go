package image

type (
	Preallocation   string
	Compat          string
	CompressionType string
	Format          string
)

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
)
