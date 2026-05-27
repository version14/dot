package shadcnui

import (
	"fmt"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	cssPath := "src/styles/globals.css"
	rsc := "false"
	if framework == "next" {
		cssPath = "src/app/globals.css"
		rsc = "true"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"tailwindcss":              "^4.0.0",
				"class-variance-authority": "^0.7.0",
				"clsx":                     "^2.0.0",
				"tailwind-merge":           "^2.0.0",
				"lucide-react":             "^0.400.0",
			},
			"devDependencies": map[string]interface{}{
				"@tailwindcss/vite": "^4.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	if err := ctx.State.UpdateJSON("tsconfig.json", func(d *state.JSONDoc) error {
		_ = d.SetNested("compilerOptions.paths", map[string]interface{}{
			"@/*": []interface{}{"./src/*"},
		})
		return nil
	}); err != nil {
		return err
	}

	componentsJSON := fmt.Sprintf(componentsJSONTmpl, rsc, cssPath)
	ctx.State.WriteFile("components.json", []byte(componentsJSON), state.ContentRaw)
	ctx.State.WriteFile(cssPath, []byte(globalCSS), state.ContentRaw)
	ctx.State.WriteFile("src/lib/utils.ts", []byte(utilsTS), state.ContentRaw)
	ctx.State.WriteFile("src/components/ui/button.tsx", []byte(buttonTSX), state.ContentRaw)

	return nil
}

const componentsJSONTmpl = `{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "default",
  "rsc": %s,
  "tsx": true,
  "tailwind": {
    "css": "%s",
    "baseColor": "slate",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils"
  }
}
`

const globalCSS = `@import "tailwindcss";
`

const utilsTS = `import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
`

const buttonTSX = `import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        outline: "border border-input hover:bg-accent hover:text-accent-foreground",
        ghost: "hover:bg-accent hover:text-accent-foreground",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 px-3",
        lg: "h-11 px-8",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => (
    <button
      className={cn(buttonVariants({ variant, size, className }))}
      ref={ref}
      {...props}
    />
  ),
);
Button.displayName = "Button";

export { Button, buttonVariants };
`
