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
		this.setGraph({
			rankdir: direction,
			nodesep: 80,
			ranksep: 120,
			marginx: 50,
			marginy: 50
		});
		this.setDefaultEdgeLabel(() => ({}));

		// Add nodes to dagre
		nodes.forEach((node) => {
			this.setNode(node.id, { width: nodeWidth, height: nodeHeight });
		});

		// Only add parent-child edges to dagre for hierarchical layout
		edges.forEach((edge) => {
			if (String(edge.data!.type).toLowerCase() === 'child') {
				this.setEdge(edge.source, edge.target);
			}
		});

		// Run dagre layout
		dagre.layout(this);

		// Create maps for relationship analysis
		const spouseMap = new Map<string, string>();
		const siblingMap = new Map<string, string[]>();
		const childrenMap = new Map<string, string[]>();
		const parentsMap = new Map<string, string[]>();

		// Build relationship maps
		edges.forEach((edge) => {
			const type = String(edge.data?.type).toLowerCase();
			if (type === 'spouse') {
				spouseMap.set(edge.source, edge.target);
				spouseMap.set(edge.target, edge.source);
			} else if (type === 'sibling') {
				if (!siblingMap.has(edge.source)) siblingMap.set(edge.source, []);
				if (!siblingMap.has(edge.target)) siblingMap.set(edge.target, []);
				siblingMap.get(edge.source)!.push(edge.target);
				siblingMap.get(edge.target)!.push(edge.source);
			} else if (type === 'child') {
				if (!childrenMap.has(edge.source)) childrenMap.set(edge.source, []);
				if (!parentsMap.has(edge.target)) parentsMap.set(edge.target, []);
				childrenMap.get(edge.source)!.push(edge.target);
				parentsMap.get(edge.target)!.push(edge.source);
			}
		});

		// Helper function to check if position is occupied
		const isPositionOccupied = (x: number, y: number, excludeId?: string): boolean => {
			return nodes.some((node) => {
				if (excludeId && node.id === excludeId) return false;
				const nodePos = this.node(node.id);
				return Math.abs(nodePos.x - x) < nodeWidth && Math.abs(nodePos.y - y) < nodeHeight;
			});
		};

		// Helper function to find free position near a reference point
		const findFreePosition = (
			refX: number,
			refY: number,
			excludeId?: string
		): { x: number; y: number } => {
			const padding = 30;
			const positions = [
				{ x: refX + nodeWidth + padding, y: refY }, // right
				{ x: refX - nodeWidth - padding, y: refY }, // left
				{ x: refX + (nodeWidth + padding) * 2, y: refY }, // far right
				{ x: refX - (nodeWidth + padding) * 2, y: refY } // far left
			];

			for (const pos of positions) {
				if (!isPositionOccupied(pos.x, pos.y, excludeId)) {
					return pos;
				}
			}

			// If all positions are taken, just go further right
			return { x: refX + (nodeWidth + padding) * 3, y: refY };
		};

		// Position spouses next to each other
		const processedSpouses = new Set<string>();
		spouseMap.forEach((spouse, person) => {
			if (processedSpouses.has(person)) return;

			const personNode = this.node(person);
			const spouseNode = this.node(spouse);

			// Determine who should be the anchor (prefer the one with children or hierarchically positioned)
			const personHasChildren = childrenMap.has(person) && childrenMap.get(person)!.length > 0;
			const spouseHasChildren = childrenMap.has(spouse) && childrenMap.get(spouse)!.length > 0;

			let anchorNode, mobileNode, anchorId, mobileId;
			if (personHasChildren && !spouseHasChildren) {
				anchorNode = personNode;
				mobileNode = spouseNode;
				anchorId = person;
				mobileId = spouse;
			} else if (!personHasChildren && spouseHasChildren) {
				anchorNode = spouseNode;
				mobileNode = personNode;
				anchorId = spouse;
				mobileId = person;
			} else {
				// Both or neither have children, use alphabetical order
				if (person < spouse) {
					anchorNode = personNode;
					mobileNode = spouseNode;
					anchorId = person;
					mobileId = spouse;
				} else {
					anchorNode = spouseNode;
					mobileNode = personNode;
					anchorId = spouse;
					mobileId = person;
				}
			}

			// Position mobile spouse next to anchor
			const newPos = findFreePosition(anchorNode.x, anchorNode.y, mobileId);
			mobileNode.x = newPos.x;
			mobileNode.y = newPos.y;

			processedSpouses.add(person);
			processedSpouses.add(spouse);
		});

		// Position siblings
		const processedSiblings = new Set<string>();
		siblingMap.forEach((siblings, person) => {
			if (processedSiblings.has(person)) return;

			const personNode = this.node(person);

			siblings.forEach((sibling) => {
				if (processedSiblings.has(sibling)) return;

				const siblingNode = this.node(sibling);

				// If sibling is not a spouse of someone at the same level, position as sibling
				const siblingSpouse = spouseMap.get(sibling);
				if (!siblingSpouse || Math.abs(this.node(siblingSpouse).y - personNode.y) > nodeHeight) {
					const newPos = findFreePosition(personNode.x, personNode.y, sibling);
					siblingNode.x = newPos.x;
					siblingNode.y = newPos.y;
				}

				processedSiblings.add(sibling);
			});

			processedSiblings.add(person);
		});

		// Create new edges
		let newEdges: Edge[] = [];
		const processedSpouseEdges = new Set<string>();

		edges.forEach((edge) => {
			let newEdge = { ...edge };

			if (String(edge.data?.type).toLowerCase() === 'child') {
				// Parent to child: source (parent) uses 'child' handle (bottom), target (child) uses 'parent' handle (top)
				newEdge.sourceHandle = 'child';
				newEdge.targetHandle = 'parent';
			} else if (String(edge.data?.type).toLowerCase() === 'parent') {
				return;
			} else if (String(edge.data?.type).toLowerCase() === 'spouse') {
				// Avoid duplicate spouse edges by creating a unique key
				const spouseKey = [edge.source, edge.target].sort().join('-');
				if (processedSpouseEdges.has(spouseKey)) {
					return; // Skip this duplicate spouse edge
				}
				processedSpouseEdges.add(spouseKey);

				// Set spouse handles based on position
				const sourceNode = this.node(edge.source);
				const targetNode = this.node(edge.target);
				if (sourceNode.x < targetNode.x) {
					newEdge.sourceHandle = 'spouse-right';
					newEdge.targetHandle = 'spouse-left';
				} else {
					newEdge.sourceHandle = 'spouse-left';
					newEdge.targetHandle = 'spouse-right';
				}
			} else if (String(edge.data?.type).toLowerCase() === 'sibling') {
				// Set sibling handles based on position
				const sourceNode = this.node(edge.source);
				const targetNode = this.node(edge.target);
				if (sourceNode.x < targetNode.x) {
					newEdge.sourceHandle = 'spouse-right';
					newEdge.targetHandle = 'spouse-left';
				} else {
					newEdge.sourceHandle = 'spouse-left';
					newEdge.targetHandle = 'spouse-right';
				}
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
