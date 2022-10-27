package types

type RestHandlers []*RestHandler

type Restful map[string]RestHandlers

type RestVersion struct {
	Plus  bool
	Value string
}

type RestHandler struct {
	Sort    int
	Sign    bool
	Login   bool
	Path    string
	Version RestVersion
	Http    string
	Method  string
	Https   map[string]byte
	Handler IController
}
