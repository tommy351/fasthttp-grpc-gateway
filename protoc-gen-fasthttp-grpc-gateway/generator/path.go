package generator

import "regexp"

var (
	pathPattern = regexp.MustCompile("{(.+?)}")
)

type Path struct {
	Path   string
	Params []string
}

func NewPath(str string) (*Path, error) {
	matches := pathPattern.FindAllStringSubmatch(str, -1)
	params := []string{}

	for _, match := range matches {
		params = append(params, match[1])
	}

	return &Path{
		Path:   pathPattern.ReplaceAllString(str, ":$1"),
		Params: params,
	}, nil
}
