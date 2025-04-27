import type { Node, Edge } from '@xyflow/svelte';

export type NodeMenu = {
	onClick: () => void;
	deleteNode: () => void;
	createRelationshipAndNode: () => void;
	addRelationship: () => void;
	addRecipe: (() => void) | undefined;
	addAdmin: (() => void) | undefined;
	top: number | undefined;
	left: number | undefined;
	right: number | undefined;
	bottom: number | undefined;
};

export type Layout = {
    Nodes: Array<Node>;
    Edges: Array<Edge>;
}