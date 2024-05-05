import { fetch_family_tree } from '$lib/family_tree/getFamilyTree';
import { getLayoutedElements } from '$lib/family_tree/dagreLayout';
import { AddToNodesData, AddToEdgesData } from '$lib/family_tree/dataAdapter';
import type { Node, Edge } from '@xyflow/svelte';
import { useEdges, useNodes } from '@xyflow/svelte';

const nodes = useNodes();
const edges = useEdges();

export function setFamilyTreeNodes() {
    console.log('fetching nodes');
    fetch_family_tree().then((data) => {
        let nodes_data: Node[] = [];
        function pushNodeToData(node: Node) {
            nodes_data.push(node);
        }

        AddToNodesData(data, 0, pushNodeToData);
        AddToNodesData(data, 2, pushNodeToData);
        AddToNodesData(data, 4, pushNodeToData);
        AddToNodesData(data, 6, pushNodeToData);
        AddToNodesData(data, 8, pushNodeToData);

        let edges_data: Edge[] = [];
        function pushEdgeToData(edge: Edge) {
            edges_data.push(edge);
        }

        AddToEdgesData(data, 1, pushEdgeToData);
        AddToEdgesData(data, 3, pushEdgeToData);
        AddToEdgesData(data, 5, pushEdgeToData);
        AddToEdgesData(data, 7, pushEdgeToData);

        const layoutedElements = getLayoutedElements(nodes_data, edges_data, 'TB');

        $nodes = layoutedElements.nodes;
        $edges = layoutedElements.edges;
    });
    console.log('nodes fetched and set');
}
