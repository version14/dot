package prettierfrontendrules

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "prettier_frontend_rules",
	Version:     "0.1.0",
	Description: "Frontend-specific Prettier rules for JSX/TSX projects",
	DependsOn:   []string{"prettier_config"},
	Outputs:     []string{},
	Validators:  []dotapi.Validator{},
}
