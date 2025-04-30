import type { Edge } from '@xyflow/svelte';

export interface RelationshipMenu {
	edge: Edge;
	XUserId: string;
	top: number | undefined;
	left: number | undefined;
	right: number | undefined;
	bottom: number | undefined;
	onClick: () => void;
	deleteEdge: () => void;
}