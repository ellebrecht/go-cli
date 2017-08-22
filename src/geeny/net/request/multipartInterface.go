package request

// MultipartInterface defines data needed for a multipart http request
type MultipartInterface interface {
	FilePath() string
	FileName() string
}
