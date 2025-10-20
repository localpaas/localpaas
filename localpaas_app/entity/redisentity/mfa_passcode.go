package redisentity

type MFAPasscode struct {
	Secret   string `json:"secret"`
	Attempts int    `json:"attempts,omitempty"`
}
