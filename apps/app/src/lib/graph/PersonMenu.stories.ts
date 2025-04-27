import type { Meta, StoryObj } from '@storybook/svelte';
import PersonMenu from './PersonMenu.svelte';

const meta = {
	title: 'graph/PersonMenu',
	component: PersonMenu,
	tags: ['autodocs'],
	argTypes: {
		top: { control: { type: 'number' } },
		left: { control: { type: 'number' } },
		right: { control: { type: 'number' } },
		bottom: { control: { type: 'number' } },
		onClick: { action: 'clicked' },
		deleteNode: { action: 'deleteNode clicked' },
		createRelationshipAndNode: { action: 'createRelationshipAndNode clicked' },
		addRelationship: { action: 'addRelationship clicked' },
		addAdmin: { action: 'addAdmin clicked' },
	},
} satisfies Meta<PersonMenu>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		top: 100,
		left: 100,
		right: undefined,
		bottom: undefined,
		onClick: () => console.log('menu closed'),
		deleteNode: () => console.log('delete node'),
		createRelationshipAndNode: () => console.log('create relationship and node'),
		addRelationship: () => console.log('add relationship'),
		addAdmin: () => console.log('add admin'),
	}
};

/**
 * Simulate a click roughly in top-left of the screen.
 */
export const TopLeftPosition: Story = {
	args: {
		top: 150,
		left: 200,
		right: undefined,
		bottom: undefined,
		onClick: () => console.log('menu closed'),
		deleteNode: () => console.log('delete node'),
		createRelationshipAndNode: () => console.log('create relationship and node'),
		addRelationship: () => console.log('add relationship'),
		addAdmin: () => console.log('add admin'),
	}
};

/**
 * Simulate a click near the bottom-right of the screen.
 */
export const BottomRightPosition: Story = {
	args: {
		top: undefined,
		left: undefined,
		right: 100,
		bottom: 120,
		onClick: () => console.log('menu closed'),
		deleteNode: () => console.log('delete node'),
		createRelationshipAndNode: () => console.log('create relationship and node'),
		addRelationship: () => console.log('add relationship'),
		addAdmin: () => console.log('add admin'),
	}
};

/**
 * Simulate a click near the middle-right (right edge but middle vertically).
 */
export const MiddleRightPosition: Story = {
	args: {
		top: 400,
		left: undefined,
		right: 50,
		bottom: undefined,
		onClick: () => console.log('menu closed'),
		deleteNode: () => console.log('delete node'),
		createRelationshipAndNode: () => console.log('create relationship and node'),
		addRelationship: () => console.log('add relationship'),
		addAdmin: () => console.log('add admin'),
	}
};
