<script>
	import PersonNode from './../lib/family_tree/PersonNode.svelte';
	import { writable } from 'svelte/store';
	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';

	const nodes = writable([
		{
			id: '1',
			data: { label: 'Hello' },
			position: { x: 0, y: 0 }
		},
		{
			id: '2',
			data: { label: 'World' },
			position: { x: 0, y: 150 }
		},
		{
			id: '3',
			data: { label: 'World' },
			position: { x: 250, y: 150 },
			type: 'custom'
		}
	]);
	const edges = writable([
		{
			id: '1-2',
			source: '1',
			target: '2'
		},
		{
			id: '1-3',
			source: '1',
			target: '3'
		}
	]);
	const nodeTypes = {
		custom: PersonNode
	};
</script>

<div style="height:100vh;">
	<SvelteFlowProvider>
		<SvelteFlow {nodes} {nodeTypes} {edges} class="bg-base-100" fitView onlyRenderVisibleElements>
			<Controls />
			<MiniMap />
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
