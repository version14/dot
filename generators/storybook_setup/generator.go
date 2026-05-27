package storybooksetup

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

	frameworkPkg := "@storybook/react-vite"
	if framework == "next" {
		frameworkPkg = "@storybook/nextjs"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"storybook":       "storybook dev -p 6006",
				"build-storybook": "storybook build",
			},
			"devDependencies": map[string]interface{}{
				"storybook":        "^10.0.0",
				"@storybook/react": "^10.0.0",
				frameworkPkg:       "^10.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	mainTS := fmt.Sprintf(mainTSTmpl, frameworkPkg, frameworkPkg)
	ctx.State.WriteFile(".storybook/main.ts", []byte(mainTS), state.ContentRaw)
	ctx.State.WriteFile(".storybook/preview.ts", []byte(previewTS), state.ContentRaw)
	ctx.State.WriteFile("src/stories/Button.stories.tsx", []byte(buttonStories), state.ContentRaw)

	return nil
}

const mainTSTmpl = `import type { StorybookConfig } from "%s";

const config: StorybookConfig = {
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|mjs|ts|tsx)"],
  addons: [],
  framework: { name: "%s", options: {} },
};

export default config;
`

const previewTS = `import type { Preview } from "@storybook/react";

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
};

export default preview;
`

const buttonStories = `import type { Meta, StoryObj } from "@storybook/react";

const ButtonComponent = ({ label }: { label: string }) => (
  <button type="button">{label}</button>
);

const meta = {
  title: "Example/Button",
  component: ButtonComponent,
  tags: ["autodocs"],
} satisfies Meta<typeof ButtonComponent>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: { label: "Button" },
};
`
