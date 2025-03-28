package memory

/* Configure guest RAM
Note: Some architectures might enforce a specific granularity
*/
type Memory struct {
	// Initial amount of guest memory
	Size uint `json:"size,omitempty"`
	// Number of hotplug slots (default: none)
	Slots uint `json:"slots,omitempty"`
	// Maximum amount of guest memory (default: none)
	MaxMem uint `json:"maxmem,omitempty"`
}
