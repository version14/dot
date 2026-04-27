package state

// ContentType tags how a FileNode's bytes should be interpreted by the
// virtual filesystem's structured update operations.
type ContentType int

const (
	ContentRaw ContentType = iota
	ContentJSON
	ContentYAML
	ContentGoMod
)

func (c ContentType) String() string {
	switch c {
	case ContentJSON:
		return "json"
	case ContentYAML:
		return "yaml"
	case ContentGoMod:
		return "gomod"
	default:
		return "raw"
	}
}
