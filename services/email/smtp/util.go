package smtp

import (
	"gopkg.in/gomail.v2"
)

var (
	msg = gomail.NewMessage()
)

// formatAddress format address as of form: "John Doe" <john.doe@address>
func formatAddress(address, displayName string) string {
	return msg.FormatAddress(address, displayName)
}
