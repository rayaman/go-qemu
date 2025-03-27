package types

/*
Note: Different machines may have different subsets of the CPU topology
      parameters supported, so the actual meaning of the supported parameters
      will vary accordingly. For example, for a machine type that supports a
      three-level CPU hierarchy of sockets/cores/threads, the parameters will
      sequentially mean as below:
                sockets means the number of sockets on the machine board
                cores means the number of cores in one socket
                threads means the number of threads in one core
      For a particular machine type board, an expected CPU topology hierarchy
      can be defined through the supported sub-option. Unsupported parameters
      can also be provided in addition to the sub-option, but their values
      must be set as 1 in the purpose of correct parsing.
*/
type SMP struct {
	// The number of initial CPUs
	Cpus uint `json:"cpus,omitempty"`
	// Maximum number of total CPUs, including offline CPUs for hotplug, etc
	MaxCpus uint `json:"maxcpus,omitempty"`
	// Number of drawers on the machine board
	Drawers uint `json:"drawers,omitempty"`
	// Number of books in one drawer
	Books uint `json:"books,omitempty"`
	// Number of sockets in one book
	Sockets uint `json:"sockets,omitempty"`
	// Number of dies in one socket
	Dies uint `json:"dies,omitempty"`
	// Number of clusters in one die
	Clusters uint `json:"clusters,omitempty"`
	// Number of modules in one cluster
	Modules uint `json:"modules,omitempty"`
	// Number of cores in one module
	Cores uint `json:"cores,omitempty"`
	// Number of threads in one core
	Threads uint `json:"threads,omitempty"`
}
