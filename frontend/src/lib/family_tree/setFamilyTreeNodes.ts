import { fetch_family_tree } from '$lib/family_tree/getFamilyTree';
import { getLayoutedElements } from '$lib/family_tree/dagreLayout';
import { AddToNodesData, AddToEdgesData } from '$lib/family_tree/dataAdapter';
import type { Node, Edge } from '@xyflow/svelte';
import { useEdges, useNodes } from '@xyflow/svelte';
import { fetch_profile } from './getProfile';

export async function setFamilyTreeNodes(): Promise<{
	nodes: Node[];
	edges: Edge[];
}> {
	console.log('fetching nodes');
	let layoutedElements: {
		nodes: Node[];
		edges: Edge[];
	} = { nodes: [], edges: [] };
	let nodes_data: Node[] = [];
	let edges_data: Edge[] = [];
	await fetch_profile().then((data) => {
		nodes_data.push({
			id: data.ElementId,
			type: 'custom',
			data: data.Props,
			position: { x: 0, y: 0 }
		});
	});

	await fetch_family_tree().then((data: []) => {
		if (data.length == 0) {
			return layoutedElements;
		}

		function pushNodeToData(node: Node) {
			nodes_data.push(node);
		}

		AddToNodesData(data, 0, pushNodeToData);
		AddToNodesData(data, 2, pushNodeToData);
		AddToNodesData(data, 4, pushNodeToData);
		AddToNodesData(data, 6, pushNodeToData);
		AddToNodesData(data, 8, pushNodeToData);

		function pushEdgeToData(edge: Edge) {
			edges_data.push(edge);
		}

		AddToEdgesData(data, 1, pushEdgeToData);
		AddToEdgesData(data, 3, pushEdgeToData);
		AddToEdgesData(data, 5, pushEdgeToData);
		AddToEdgesData(data, 7, pushEdgeToData);
	});

	console.log('Fetched nodes and edges data.');

	// Remove duplicate nodes
	nodes_data = nodes_data.filter(
		(node, index, self) => index === self.findIndex((n) => n.id === node.id)
	);

	return getLayoutedElements(nodes_data, edges_data, 'TB');
}
