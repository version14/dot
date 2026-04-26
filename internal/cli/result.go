package cli


// AnswerEntry is one answered question in flow order.
type AnswerEntry struct {
    Key        string                 `json:"key"`
    Value      interface{}            `json:"value,omitempty"`      // string | bool | int
    Multi      []string               `json:"multi,omitempty"`      // multi-select
    Iterations []map[string]interface{} `json:"iterations,omitempty"` // loop iterations
}

// Result holds answers in the order the user encountered them.
type Result struct {
    Entries []AnswerEntry
    index   map[string]int // key → position for O(1) lookup
}

// Typed accessors — avoid interface{} assertions in generator code
func (r *Result) GetString(key string) string { ... }
func (r *Result) GetBool(key string) bool     { ... }
func (r *Result) GetInt(key string) int       { ... }
func (r *Result) GetMulti(key string) []string { ... }
