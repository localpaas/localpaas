package base

type FileCompressionFormat string

const (
	FileCompressionNone       FileCompressionFormat = ""
	FileCompressionFormatGzip FileCompressionFormat = "gzip"
)

var (
	AllFileCompressionFormats = []FileCompressionFormat{FileCompressionNone, FileCompressionFormatGzip}
)
