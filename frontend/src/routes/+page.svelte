<script lang="ts">
	import { writable } from 'svelte/store';
	import { onMount } from 'svelte';
	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap, Background } from '@xyflow/svelte';
	import type { Node, Edge, NodeTypes } from '@xyflow/svelte';
	import { isAuthenticated } from '../lib/stores';
	import PersonNode from './../lib/family_tree/PersonNode.svelte';
	import { login } from '../lib/auth';
	import { fetch_family_tree } from '$lib/family_tree/getFamilyTree';
	import { getLayoutedElements } from '$lib/family_tree/dagreLayout';
	import { AddToNodesData, AddToEdgesData } from '$lib/family_tree/dataAdapter';
	let edges_data: {
		id: string;
		source: string;
		target: string;
		data: { type: string; verified: boolean };
	}[] = [];
	const nodes = writable<Node[]>([]);
	const edges = writable<Edge[]>([]);
	onMount(() => {
		if (!$isAuthenticated) {
			console.log('user authenticated:', $isAuthenticated);
			login();
		} else {
			console.log('user authenticated:', $isAuthenticated);
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
		}
	});

	const nodeTypes: NodeTypes = {
		custom: PersonNode
	};

	let menu: { id: string; top?: number; left?: number; right?: number; bottom?: number } | null;
	let width: number;
	let height: number;

	function handleContextMenu({ detail: { event, node } }) {
		// Prevent native context menu from showing
		event.preventDefault();

		// Calculate position of the context menu. We want to make sure it
		// doesn't get positioned off-screen.
		menu = {
			id: node.id,
			top: event.clientY < height - 200 ? event.clientY : undefined,
			left: event.clientX < width - 200 ? event.clientX : undefined,
			right: event.clientX >= width - 200 ? width - event.clientX : undefined,
			bottom: event.clientY >= height - 200 ? height - event.clientY : undefined
		};
	}

	// Close the context menu if it's open whenever the window is clicked.
	function handlePaneClick() {
		menu = null;
	}
</script>

<div style="height:100vh;">
	<SvelteFlowProvider>
		<SvelteFlow {nodes} {nodeTypes} {edges} class="bg-base-100" fitView onlyRenderVisibleElements>
			<Controls class="bg-base-300 text-base-content" />
			<MiniMap class="bg-base-200"/>
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
