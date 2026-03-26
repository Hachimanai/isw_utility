package domain

import "time"

// SystemInfo represents Kernel version, Uptime, and CPU frequency.
type SystemInfo struct {
	KernelVersion string        `json:"kernel_version"`
	Uptime        time.Duration `json:"uptime"`
	CPUFreq       float64       `json:"cpu_freq"` // Frequency in GHz or MHz
}
