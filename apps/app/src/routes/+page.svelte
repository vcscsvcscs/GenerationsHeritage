<script lang="ts">
	import { onMount } from 'svelte';
	import { nodeTypes, edgeTypes } from '$lib/graph/model';
	import { title, family_tree } from '$lib/paraglide/messages.js';

	import {
		SvelteFlowProvider,
		SvelteFlow,
		Controls,
		MiniMap
	} from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import type { OnConnectEnd, Node, Edge, NodeEventWithPointer } from '@xyflow/svelte';

	import PersonModal from '$lib/profile/Modal.svelte';
	import PersonMenu from '$lib/graph/PersonMenu.svelte';
	import CreatePerson from '$lib/profile/create/Modal.svelte';

	import type { components } from '$lib/api/api.gen';
	import type { NodeMenu } from '$lib/graph/model';

	import { handleNodeClick } from '$lib/graph/node_click';

	import { FamilyTree } from '$lib/graph/layout';
	import { tailwindClassToPixels } from '$lib/tailwindSizeToPx';
	import type { Layout } from '$lib/graph/model';
	import HamburgerIcon from '$lib/sidebar/hamburgerIcon.svelte';

	let { data }: { data: Layout & { id: string } } = $props();

	let selectedPerson: components['schemas']['PersonProperties'] & { id: number | null } = $state({
		id: null
	});
	let openPersonPanel = $state(false);
	let openPersonMenu: NodeMenu | undefined = $state(undefined);
	let with_out_spouse = $state(false);

	let familyTreeDAG = new FamilyTree();
	let layout = familyTreeDAG.getLayoutedElements(
		data.Nodes,
		data.Edges,
		tailwindClassToPixels('w-40') || 160,
		tailwindClassToPixels('h-40') || 160,
		'TB'
	);
	console.log('layout', layout);
	let nodes = $state.raw<Node[]>([] as Node[]);
	let edges = $state.raw<Edge[]>([] as Edge[]);

	let relationshipStart: number | null = $state(null);
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
			XUserId: data.id,
			onClick: () => {
				openPersonMenu = undefined;
			},
			deleteNode: () => {
				if (Number(data.id) === Number(node.data.id)) {
					relationshipStart = null;
					openPersonMenu = undefined;

					return;
				}
				fetch('/api/person/' + node.data.id, {
					method: 'DELETE',
					headers: {
						'Content-Type': 'application/json'
					}
				})
					.then((response) => {
						if (response.ok) {
							nodes = nodes.filter((n) => n.data.id !== node.data.id);
							edges = edges.filter((e) => e.source !== "help"+node.data.id && e.target !== "help"+node.data.id);
						} else {
							alert('Error deleting person');
						}
					})
					.catch((error) => {
						console.error('Error:', error);
					});
				openPersonMenu = undefined;
			},
			createRelationshipAndNode: () => {
				relationshipStart = Number(node.data.id);
				createPerson = true;
				openPersonMenu = undefined;
			},
			addRelationship: () => {
				relationshipStart = Number(node.data.id);
				openPersonMenu = undefined;
			},
			addAdmin: () => {
				relationshipStart = Number(node.data.id);
				openPersonMenu = undefined;
			},
			addRecipe: () => {
				relationshipStart = Number(node.data.id);
				openPersonMenu = undefined;
			},
			id: String(node.data.id),
			top: event.clientY < clientHeight - 200 ? event.clientY : undefined,
			left: event.clientX < clientWidth - 200 ? event.clientX : undefined,
			right: event.clientX >= clientWidth - 200 ? clientWidth - event.clientX : undefined,
			bottom: event.clientY >= clientHeight - 200 ? clientHeight - event.clientY : undefined
		};
	};

	let onCreation = (newNodes: Array<Node> | null, newEdges: Array<Edge> | null) => {
		if (newNodes !== null) {
			nodes = [...nodes, ...newNodes];
		}

		if (newEdges !== null) {
			edges = [...edges, ...newEdges];
		}

		console.log('newnodes', newNodes![0].id);
		console.log('newedges', newEdges![0].id);
		console.log('newnodes', newNodes);
		console.log('newedges', newEdges);

		let newLayout = familyTreeDAG.getLayoutedElements(
			nodes,
			edges,
			tailwindClassToPixels('w-40') || 160,
			tailwindClassToPixels('h-40') || 160,
			'TB'
		);
		edges = [...newLayout.Edges];
		nodes = [...newLayout.Nodes];
	};

	let handleNodeClickFunc = handleNodeClick(
		(
			person: components['schemas']['PersonProperties'] & {
				id: number;
			}
		) => {
			openPersonPanel = true;
			selectedPerson = person;
			fetch('/api/person/' + person.id, {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json'
				}
			})
				.then((response) => {
					if (response.ok) {
						return response.json() as Promise<components['schemas']['Person']>;
					} else {
						alert('Error fetching person data');
						return null;
					}
				})
				.then((data) => {
					if (data) {
						selectedPerson = data.Props as components['schemas']['PersonProperties'] & {
							id: number | null;
						};
						selectedPerson.id = person.id;
					}
				});
		}
	);

	let handlePaneClick = ({ event }: { event: MouseEvent }) => {
		openPersonPanel = false;
		openPersonMenu = undefined;
	};

	const handleConnectEnd: OnConnectEnd = (event, connectionState) => {
		if (connectionState.isValid) return;
		const sourceNodeId = connectionState.fromNode?.data.id;
		if (sourceNodeId === undefined) return;
		relationshipStart = Number(sourceNodeId);
		createPerson = true;
	};
	onMount(() => {
		nodes = [...layout.Nodes];
		edges = [...layout.Edges]
	});
</script>

<svelte:head>
	<title>{title({ page: family_tree() })}</title>
</svelte:head>
<div style="height:100vh;" class="!bg-base-200 flex flex-col">
	<SvelteFlowProvider>
		<SvelteFlow
			bind:nodes={nodes}
			bind:edges={edges}
			onconnectend={handleConnectEnd}
			onnodeclick={handleNodeClickFunc}
			onnodecontextmenu={handleContextMenu}
			onpaneclick={handlePaneClick}
			class="!bg-base-200"
			{nodeTypes}
			{edgeTypes}
			fitView
			onlyRenderVisibleElements={false}
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
				<CreatePerson
					onOnlyPersonCreation={() => {
						createPerson = false;
					}}
					{onCreation}
					closeModal={() => {
						createPerson = false;
					}}
					relationshipStartID={relationshipStart}
				></CreatePerson>
			{/if}
			{#if openPersonMenu !== undefined}
				<PersonMenu {...openPersonMenu!} />
			{/if}
		</SvelteFlow>
	</SvelteFlowProvider>
</div>

<div class="absolute top-2 left-2 flex flex-row items-center gap-2">
	<HamburgerIcon />
</div>
