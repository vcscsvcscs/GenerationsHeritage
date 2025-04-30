import dagre from '@dagrejs/dagre';
import type { Layout } from './model';
import type { Edge, Node } from '@xyflow/svelte';
import { Position } from '@xyflow/svelte';

export class FamilyTree extends dagre.graphlib.Graph {
	constructor() {
		super();
	}

	getLayoutedElements(
		nodes: Node[],
		edges: Edge[],
		nodeWidth: number,
		nodeHeight: number,
		direction = 'TB'
	): Layout {
		this.setGraph({ rankdir: direction });
		this.setDefaultEdgeLabel(() => ({}));
		nodes.forEach((node) => {
			this.setNode(node.id, { width: nodeWidth, height: nodeHeight });
		});

		edges.forEach((edge) => {
			if (String(edge.data!.type).toLowerCase() === 'child') {
				this.setEdge(edge.source, edge.target);
			}
		});

		dagre.layout(this);

		let newEdges: Edge[] = [];
		edges.forEach((edge) => {
			let newEdge = { ...edge };
			if (String(edge.data?.type).toLowerCase() === 'child') {
				newEdge.sourceHandle = 'child';
				newEdge.targetHandle = 'parent';
			}
			
			const sourceNode = this.node(edge.source);
			const targetNode = this.node(edge.target);
			if (!sourceNode || !targetNode) {
				return;
			}

			if (String(edge.data?.type).toLowerCase() === 'spouse') {
				const padding = 50; // distance between spouse and source
				const spouseWidth = nodeWidth;

				const existingNodesAtLevel = nodes
					.map((n) => ({ id: n.id, pos: this.node(n.id) }))
					.filter(({ pos }) => Math.abs(pos.y - sourceNode.y) < nodeHeight / 2); // same horizontal band

				// Collect taken x ranges
				const takenXRanges = existingNodesAtLevel.map(({ pos }) => ({
					from: pos.x - spouseWidth / 2,
					to: pos.x + spouseWidth / 2
				}));

				// Try placing spouse to the right
				let desiredX = sourceNode.x + nodeWidth + padding;

				// Check for collision
				const collides = (x: number) => {
					return takenXRanges.some(({ from, to }) => x > from && x < to);
				};

				// If right side collides, try left
				if (collides(desiredX)) {
					desiredX = sourceNode.x - (nodeWidth + padding);
				}

				// If both sides collide, push right until free
				while (collides(desiredX)) {
					desiredX += nodeWidth + padding;
				}

				targetNode.x = desiredX;
				targetNode.y = sourceNode.y;
			}
			newEdge.hidden = false;
			newEdge.type = 'familyEdge';

			newEdges.push(newEdge);
		});

		const layoutedNodes = nodes.map((node) => {
			const nodeWithPosition = this.node(node.id);

			return {
				...node,
				type: 'personNode',
				position: {
					x: nodeWithPosition.x - nodeWidth / 2,
					y: nodeWithPosition.y - nodeHeight / 2
				}
			};
		});

		return { Nodes: layoutedNodes, Edges: newEdges };
	}
}
