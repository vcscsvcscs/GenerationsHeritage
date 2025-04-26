<script lang="ts">
	import { title, family_tree, people } from '$lib/paraglide/messages.js';
	import type { PageProps } from './$types';
	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import type { Node, Edge, NodeTypes, NodeProps } from '@xyflow/svelte';
	import PersonNode from '$lib/graph/PersonNode.svelte';
	import PersonModal from '$lib/profile/Modal.svelte';
	import type { components } from '$lib/api/api.gen';
	import { handleNodeClick } from '$lib/graph/node_click';
	import PersonMenu from '$lib/graph/PersonMenu.svelte';
	import type { NodeMenu } from '$lib/graph/model';

	let { data, form }: PageProps = $props();
	const nodeTypes: NodeTypes = { personNode: PersonNode };
	let selectedPerson: components['schemas']['PersonProperties'] & {id:number|null} = $state({id: null});
	let openPersonPanel = $state(false);
	let openPersonMenu: NodeMenu = $state({});

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

	handleNodeClick(
		(
			person: components['schemas']['PersonProperties'] & {
				id: number;
			}
		) => {
			openPersonPanel = true;
			selectedPerson = person;
		}
	);
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
			<PersonModal person={selectedPerson} open={openPersonPanel} />
			<PersonMenu />
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
