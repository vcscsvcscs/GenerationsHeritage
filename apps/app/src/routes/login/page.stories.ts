import type { Meta, StoryObj } from '@storybook/svelte';
import SignInPage from './+page.svelte';

const meta = {
	title: 'login/+page',
	component: SignInPage,
	tags: ['autodocs']
} satisfies Meta<typeof SignInPage>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => {
		return {
			Component: SignInPage
		};
	}
};
