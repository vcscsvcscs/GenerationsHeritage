import type { Node, Edge, NodeTypes, EdgeTypes } from '@xyflow/svelte';
import FamilyEdge from './FamilyEdge.svelte';
import PersonNode from './PersonNode.svelte';

export const nodeTypes: NodeTypes = { personNode: PersonNode };
export const edgeTypes: EdgeTypes = {
		familyEdge: FamilyEdge
};

export type NodeMenu = {
	onClick: () => void;
	deleteNode: () => void;
	createRelationshipAndNode: () => void;
	addRelationship: () => void;
	addRecipe: (() => void) | undefined;
	addAdmin: (() => void) | undefined;
	id: string;
	XUserId: string;
	top: number | undefined;
	left: number | undefined;
	right: number | undefined;
	bottom: number | undefined;
};

export type Layout = {
	Nodes: Array<Node>;
	Edges: Array<Edge>;
};
