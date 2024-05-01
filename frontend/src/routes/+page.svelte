<script lang="ts">
	import { writable } from 'svelte/store';
	import { onMount } from 'svelte';
	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';
	import type { Node, Edge, NodeTypes } from '@xyflow/svelte';
	import { isAuthenticated } from '../lib/stores';
	import PersonNode from './../lib/family_tree/PersonNode.svelte';
	import { login } from '../lib/auth';
	import { fetch_family_tree } from '$lib/family_tree/getFamilyTree';
	import { getLayoutedElements } from '$lib/family_tree/dagreLayout';

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
				let Position = { x: 0, y: 0 };
				let nodes_data: Node[] = [];
				function AddToNodesData(data: any, i: number) {
					if (data[0].Values[i] != null) {
						if (Object.prototype.toString.call(data[0].Values[i]) === '[object Array]') {
							data[0].Values[i].forEach((person: { ElementId: string; Props: {} }) => {
								nodes_data.push({
									id: person.ElementId,
									type: 'custom',
									data: person.Props,
									position: Position
								});
							});
						} else {
							nodes_data.push({
								id: data[0].Values[i].ElementId,
								type: 'custom',
								data: data[0].Values[i].Props,
								position: Position
							});
						}
					}
				}
				AddToNodesData(data, 0);
				AddToNodesData(data, 2);
				AddToNodesData(data, 4);
				AddToNodesData(data, 6);
				AddToNodesData(data, 8);

				edges_data = [];

				const layoutedElements = getLayoutedElements(nodes_data, edges_data, 'TB');

				$nodes = layoutedElements.nodes;
				$edges = layoutedElements.edges;
			});
		}
	});

	const nodeTypes: NodeTypes = {
		custom: PersonNode
	};
</script>

<div style="height:100vh;">
	<SvelteFlowProvider>
		<SvelteFlow {nodes} {nodeTypes} {edges} class="bg-base-100" fitView onlyRenderVisibleElements>
			<Controls class="bg-base-300 text-primary-content" />
			<MiniMap />
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
