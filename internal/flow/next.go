package flow

// Next is the edge in the flow graph. Exactly one of Question, Fragment, or End
// should be set per edge. End=true terminates traversal.
type Next struct {
	Question Question // go directly to this question node
	Fragment string   // resolve a named fragment (may inject plugin questions)
	End      bool     // end the flow
}
