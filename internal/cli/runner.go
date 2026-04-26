package cli

// Runner — pre-walk, single Huh form, sub-runners for loops


type Runner struct {
    Flow         *q.Question
    Hooks        map[string]Hook
    Result       *Result
    parentResult *Result             // outer scope answers — used by loop sub-runners
    strPtrs      map[string]*string  // live pointers written by huh during Run()
    multiPtrs    map[string]*[]string
}

func (runner *Runner) condsMatch(conditions []conditionType) bool {
    for _, condition := range conditions {
        var actual string
        if ptr, ok := runner.strPtrs[condition.key]; ok {
            actual = *ptr  // live — read during WithHideFunc evaluation
        } else {
            actual = runner.Result.Get(condition.key)
            if actual == "" && runner.parentResult != nil {
                actual = runner.parentResult.Get(condition.key) // cross-scope
            }
        }
        // ...
    }
}

func (runner *Runner) Run() error {
    // 1. Pre-walk
    var nodes []flatNode
    var loops []loopEntry
    flatten(runner.Flow, nil, &nodes, &loops)

    // 2. Build groups with reactive hide functions
    groups := runner.buildGroups(nodes)

    // 3. Single form — full back/forward navigation
    huh.NewForm(groups...).Run()

    // 4. Collect only visited answers
    for _, node := range nodes {
        if runner.condsMatch(node.conditions) {
            runner.Result.add(...)
        }
    }

    // 5. Hooks
    // 6. Sub-runners for loops
    for _, le := range loops {
        if runner.condsMatch(le.conditions) {
            runner.runLoop(le.loop) // runs N times, merges into Result
        }
    }
}
