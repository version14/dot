package question

type Option struct {
	Label       string
	Value       string
	Next        *Next
	GeneratorID string // ID of the generator to invoke when this option is selected
}

type OptionQuestion struct {
	Label       string
	Description string
	Value       string
	Multiple    bool // true = multi-select (e.g. databases)
	Options     []*Option
	Next        *Next // used when Multiple=true: continuation after selection
}
