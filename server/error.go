package server

// Constant Error values which can be compared to determine the type of error
const (
	ErrorNoSuchEmbeddedResource = iota
	ErrorNoSuchFileResource
)

// Error represents a more strongly typed server error for detecting
// and handling specific types of errors.
type Error int

func (e Error) Error() string {
	return errorMessages[e]
}

var errorMessages = map[Error]string{
	ErrorNoSuchEmbeddedResource: "No embedded resource with the given path exists",
	ErrorNoSuchFileResource:     "No file resource with the given path exists",
}
