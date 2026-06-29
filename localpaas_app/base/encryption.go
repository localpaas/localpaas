package base

const (
	EncryptionSaltPrefix = "lpsalt:"
)

type FileEncryptionFormat string

const (
	FileEncryptionNone      FileEncryptionFormat = ""
	FileEncryptionFormatAge FileEncryptionFormat = "age"
)

var (
	AllFileEncryptionFormats = []FileEncryptionFormat{FileEncryptionNone, FileEncryptionFormatAge}
)
