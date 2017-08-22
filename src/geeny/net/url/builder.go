package url

// Builder builds urls
type Builder struct {
	baseURL string
}

// NewBuilder creates a new URLBuilder instance
func NewBuilder(baseURL string) *Builder {
	return &Builder{
		baseURL: baseURL,
	}
}

// NewURL creates a new URL from given endpoint
func (b *Builder) NewURL(endpoint string) string {
	return b.baseURL + endpoint
}
