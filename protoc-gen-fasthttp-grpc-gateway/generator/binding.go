package generator

//go:generate stringer -type=HTTPMethod
type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	PUT
	PATCH
	DELETE
)

type Binding struct {
	Method     *Method
	Index      int
	HTTPMethod HTTPMethod
	Path       *Path
	Body       string
	BodyType   *Message
}
