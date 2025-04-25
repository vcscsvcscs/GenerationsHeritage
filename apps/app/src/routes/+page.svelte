<script lang="ts">
	import { title, family_tree } from '$lib/paraglide/messages.js';
	import type { PageData } from './$types';
	import {
		SvelteFlowProvider,
		Background,
		BackgroundVariant,
		SvelteFlow,
		Controls,
		MiniMap
	} from '@xyflow/svelte';
	import type { Node, Edge, NodeTypes, NodeProps } from '@xyflow/svelte';
	let data: PageData = $props();
	let nodes = $state.raw<Node[]>([
		{
			id: '1',
			type: 'input',
			data: { label: 'Input Node' },
			position: { x: 0, y: 0 }
		},
		{
			id: '2',
			data: { label: 'Default Node' },
			position: { x: 100, y: 100 }
		},
		{
			id: '3',
			data: { label: 'Output Node' },
			position: { x: 200, y: 200 }
		}
	]);
	let edges = $state.raw<Edge[]>([]);
</script>

<svelte:head>
	<title>{title({ page: family_tree() })}</title>
</svelte:head>

<div style="height:100vh;" class="flex flex-col bg-base-200">
	<SvelteFlowProvider>
		<SvelteFlow bind:nodes bind:edges class="bg-base-200" fitView onlyRenderVisibleElements>
			<MiniMap />
			<Controls />
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
