<script lang="ts">
	import { title, family_tree, people } from '$lib/paraglide/messages.js';
	import type { PageProps } from './$types';
	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import type { Node, Edge, NodeTypes, NodeProps } from '@xyflow/svelte';
	import PersonNode from '$lib/graph/PersonNode.svelte';
	import type { components } from '$lib/api/api.gen';
	let { data, form }: PageProps = $props();
	const nodeTypes: NodeTypes = { personNode: PersonNode };
	let ppl = data.people;
	if (ppl === undefined) {
		ppl = [];
	}
	if (ppl[0].id === undefined) {
		ppl[0].id = 0;
	}

	let nodes = $state.raw<Node[]>([
		{
			id: String(ppl[0].id),
			type: 'personNode',
			data: ppl[0] as components['schemas']['PersonProperties'],
			position: { x: 0, y: 0 }
		}
	]);
	let edges = $state.raw<Edge[]>([]);
</script>

<svelte:head>
	<title>{title({ page: family_tree() })}</title>
</svelte:head>

<div style="height:100vh;" class="!bg-base-200 flex flex-col">
	<SvelteFlowProvider>
		<SvelteFlow
			bind:nodes
			bind:edges
			class="!bg-base-200"
			{nodeTypes}
			fitView
			onlyRenderVisibleElements
		>
			<MiniMap class="!bg-base-300" />
			<Controls class="!bg-base-300" />
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
