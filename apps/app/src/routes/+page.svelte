<script lang="ts">
	import type { PageProps } from './$types';
	import { title, family_tree } from '$lib/paraglide/messages.js';

	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import type { Node, Edge, NodeTypes, NodeEventWithPointer } from '@xyflow/svelte';

	import PersonNode from '$lib/graph/PersonNode.svelte';
	import PersonModal from '$lib/profile/Modal.svelte';
	import PersonMenu from '$lib/graph/PersonMenu.svelte';
	import CreatePerson from '$lib/profile/create/Modal.svelte';

	import type { components } from '$lib/api/api.gen';
	import type { NodeMenu } from '$lib/graph/model';

	import { handleNodeClick } from '$lib/graph/node_click';

	let { data, form }: PageProps = $props();
	const nodeTypes: NodeTypes = { personNode: PersonNode };
	let selectedPerson: components['schemas']['PersonProperties'] & { id: number | null } = $state({
		id: null
	});
	let openPersonPanel = $state(false);
	let openPersonMenu: NodeMenu | undefined = $state(undefined);

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

	let relationshipStart: number | null ;
	let createPerson = $state(false);

	let clientWidth: number | undefined = $state();
	let clientHeight: number | undefined = $state();
	const handleContextMenu: NodeEventWithPointer<MouseEvent> = ({ event, node }) => {
		event.preventDefault();

		if (clientHeight === undefined || clientWidth === undefined) {
			clientHeight = window.innerHeight;
			clientWidth = window.innerWidth;
		}

		openPersonMenu = {
			onClick: () => {
				openPersonMenu = undefined;
			},
			deleteNode: () => {
				openPersonMenu = undefined;
			},
			createRelationshipAndNode: () => {
				openPersonMenu = undefined;
			},
			addRelationship: () => {
				openPersonMenu = undefined;
			},
			addAdmin: () => {
				openPersonMenu = undefined;
			},
			addRecipe: () => {
				openPersonMenu = undefined;
			},
			top: event.clientY < clientHeight - 200 ? event.clientY : undefined,
			left: event.clientX < clientWidth - 200 ? event.clientX : undefined,
			right: event.clientX >= clientWidth - 200 ? clientWidth - event.clientX : undefined,
			bottom: event.clientY >= clientHeight - 200 ? clientHeight - event.clientY : undefined
		};
	};

	let handleNodeClickFunc = handleNodeClick(
		(
			person: components['schemas']['PersonProperties'] & {
				id: number;
			}
		) => {
			openPersonPanel = true;
			console.log('person', person);
			selectedPerson = person;
		}
	);

	let handlePaneClick = ({ event }: { event: MouseEvent }) => {
		openPersonPanel = false;
		openPersonMenu = undefined;
	};
</script>

<svelte:head>
	<title>{title({ page: family_tree() })}</title>
</svelte:head>

<div style="height:100vh;" class="!bg-base-200 flex flex-col">
	<SvelteFlowProvider>
		<SvelteFlow
			bind:nodes
			bind:edges
			onnodeclick={handleNodeClickFunc}
			onnodecontextmenu={handleContextMenu}
			onpaneclick={handlePaneClick}
			class="!bg-base-200"
			{nodeTypes}
			fitView
			onlyRenderVisibleElements
		>
			<MiniMap class="!bg-base-300" />
			<Controls class="!bg-base-300" />
			{#if openPersonPanel}
				<PersonModal
					person={selectedPerson}
					closeModal={() => {
						openPersonPanel = false;
					}}
				/>
			{/if}
			{#if createPerson}
				<CreatePerson relationship={relationshipStart}></CreatePerson>
			{/if}
			{#if openPersonMenu !== undefined}
				<PersonMenu {...openPersonMenu!} />
			{/if}
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
