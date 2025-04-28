import type { Node, Edge, NodeTypes } from '@xyflow/svelte';
import PersonNode from './PersonNode.svelte';

export const nodeTypes: NodeTypes = { personNode: PersonNode };

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
};
