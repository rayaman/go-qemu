package formats

import (
	"github.com/rayaman/go-qemu/pkg/types"
)

func init() {
	types.RegisterType("qcow2", QCOW2{})
}

type QCOW2 struct {
	// Name of the image
	ImageName string `json:"image_name"`
	// Determines the qcow2 version to use. compat=0.10 uses the traditional image format that can be read by any QEMU since 0.10. compat=1.1 enables image format extensions that only QEMU 1.1 and newer understand (this is the default). Amongst others, this includes zero clusters, which allow efficient copy-on-read for sparse images
	Compat types.Compat `json:"compat,omitempty"`
	// File name of a base image
	BackingFile string `json:"backing_file,omitempty"`
	// Image format of the base image
	BackingFmt string `json:"backing_fmt,omitempty"`
	// This option configures which compression algorithm will be used for compressed clusters on the image. Note that setting this option doesn’t yet cause the image to actually receive compressed writes. It is most commonly used with the -c option of qemu-img convert, but can also be used with the compress filter driver or backup block jobs with compression enabled.
	CompressionType types.CompressionType `json:"compression_type,omitempty"`
	/*
		If this option is set to true, the image is encrypted with 128-bit AES-CBC.

		The use of encryption in qcow and qcow2 images is considered to be flawed by modern cryptography standards, suffering from a number of design problems:

		*	The AES-CBC cipher is used with predictable initialization vectors based on the sector number. This makes it vulnerable to chosen plaintext attacks which can reveal the existence of encrypted data.

		*	The user passphrase is directly used as the encryption key. A poorly chosen or short passphrase will compromise the security of the encryption.

		*	In the event of the passphrase being compromised there is no way to change the passphrase to protect data in any qcow images. The files must be cloned, using a different encryption passphrase in the new file. The original file must then be securely erased using a program like shred, though even this is ineffective with many modern storage technologies.

		*	Initialization vectors used to encrypt sectors are based on the guest virtual sector number, instead of the host physical sector. When a disk image has multiple internal snapshots this means that data in multiple physical sectors is encrypted with the same initialization vector. With the CBC mode, this opens the possibility of watermarking attacks if the attack can collect multiple sectors encrypted with the same IV and some predictable data. Having multiple qcow2 images with the same passphrase also exposes this weakness since the passphrase is directly used as the key.

		Use of qcow / qcow2 encryption is thus strongly discouraged. Users are recommended to use an alternative encryption technology such as the Linux dm-crypt / LUKS system.
	*/
	Encryption bool `json:"encryption,omitempty"`
	// Changes the qcow2 cluster size (must be between 512 and 2M). Smaller cluster sizes can improve the image file size whereas larger cluster sizes generally provide better performance.
	ClusterSize string `json:"cluster_size,omitempty"`
	// Preallocation mode (allowed values: off, metadata, falloc, full). An image with preallocated metadata is initially larger but can improve performance when the image needs to grow. falloc and full preallocations are like the same options of raw format, but sets up metadata also.
	Preallocation types.Preallocation `json:"preallocation,omitempty"`
	/*
		If this option is set to true, reference count updates are postponed with the goal of avoiding metadata I/O and improving performance. This is particularly interesting with cache=writethrough which doesn’t batch metadata updates. The tradeoff is that after a host crash, the reference count tables must be rebuilt, i.e. on the next open an (automatic) qemu-img check -r all is required, which may take some time.

		This option can only be enabled if compat=1.1 is specified.
	*/
	LazyRefcounts bool `json:"lazy_refcounts,omitempty"`
	/*
		If this option is set to true, it will turn off COW of the file. It’s only valid on btrfs, no effect on other file systems.

		Btrfs has low performance when hosting a VM image file, even more when the guest on the VM also using btrfs as file system. Turning off COW is a way to mitigate this bad performance. Generally there are two ways to turn off COW on btrfs:

		*	Disable it by mounting with nodatacow, then all newly created files will be NOCOW

		*	For an empty file, add the NOCOW file attribute. That’s what this option does.

		Note: this option is only valid to new or empty files. If there is an existing file which is COW and has data blocks already, it couldn’t be changed to NOCOW by setting nocow=on. One can issue lsattr filename to check if the NOCOW flag is set or not (Capital ‘C’ is NOCOW flag).
	*/
	NoCow bool `json:"nocow,omitempty"`
	/*
		Filename where all guest data will be stored. If this option is used, the qcow2 file will only contain the image’s metadata.

		Note: Data loss will occur if the given filename already exists when using this option with qemu-img create since qemu-img will create the data file anew, overwriting the file’s original contents. To simply update the reference to point to the given pre-existing file, use qemu-img amend.
	*/
	DataFile string `json:"data_file,omitempty"`
	/*
		If this option is set to true, QEMU will always keep the external data file consistent as a standalone read-only raw image.

		It does this by forwarding all write accesses to the qcow2 file through to the raw data file, including their offsets. Therefore, data that is visible on the qcow2 node (i.e., to the guest) at some offset is visible at the same offset in the raw data file. This results in a read-only raw image. Writes that bypass the qcow2 metadata may corrupt the qcow2 metadata because the out-of-band writes may result in the metadata falling out of sync with the raw image.

		If this option is off, QEMU will use the data file to store data in an arbitrary manner. The file’s content will not make sense without the accompanying qcow2 metadata. Where data is written will have no relation to its offset as seen by the guest, and some writes (specifically zero writes) may not be forwarded to the data file at all, but will only be handled by modifying qcow2 metadata.

		This option can only be enabled if data_file is set.
	*/
	DataFileRaw bool `json:"data_file_raw,omitempty"`
}
