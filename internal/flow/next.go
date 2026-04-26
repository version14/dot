package flow

// Next struct — the graph edge. Points to a Question, a Fragment, or End=true

type Next struct {
	Question *Question // go directly to this question node
	Fragment string    // resolve a named fragment (may inject plugin questions)
	End      bool      // end the flow
}
