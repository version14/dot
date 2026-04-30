package prettierexpressrules

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "prettier_express_rules",
	Version:     "0.1.0",
	Description: "Express/Node backend-specific Prettier rules (printWidth, endOfLine)",
	DependsOn:   []string{"prettier_config"},
	Outputs:     []string{},
	Validators:  []dotapi.Validator{},
}
