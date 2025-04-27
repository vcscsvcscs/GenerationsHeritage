import type { Meta, StoryObj } from '@storybook/svelte';
import Logout from './Logout.svelte';

const meta = {
	title: 'lib/Logout',
	component: Logout,
	tags: ['autodocs'],
	argTypes: {
		show: { control: { type: 'boolean' } }
	}
} satisfies Meta<Logout>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Visible: Story = {
	args: {
		show: true
	}
};

export const Hidden: Story = {
	args: {
		show: false
	}
};
