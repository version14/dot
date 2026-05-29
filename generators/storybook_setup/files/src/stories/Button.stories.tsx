import type { Meta, StoryObj } from "@storybook/react";

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
