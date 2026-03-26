package domain

// FanStatus represents Fan RPM and its Idle/Max range.
type FanStatus struct {
	Label   string `json:"label"`    // e.g., "CPU Fan", "GPU Fan"
	RPM     int    `json:"rpm"`      // Current revolutions per minute
	IdleRPM int    `json:"idle_rpm"` // Minimum RPM in idle state
	MaxRPM  int    `json:"max_rpm"`  // Maximum RPM possible
}
