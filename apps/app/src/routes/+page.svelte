<script lang="ts">
	import { nodeTypes, edgeTypes } from '$lib/graph/model';
	import { title, family_tree } from '$lib/paraglide/messages.js';

	import {
		SvelteFlowProvider,
		SvelteFlow,
		Controls,
		MiniMap,
		ConnectionLineType
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
	import SideBar from '$lib/sidebar/sideBar.svelte';

	let { data }: { data: Layout & { id: string } } = $props();

	let selectedPerson: components['schemas']['PersonProperties'] & { id: number | null } = $state({
		id: null
	});
	let openPersonPanel = $state(false);
	let openPersonMenu: NodeMenu | undefined = $state(undefined);
	let with_out_spouse = $state(false);

	let familyTreeDAG = new FamilyTree();
	let nodes = $state.raw<Node[]>([]);
	let edges = $state.raw<Edge[]>([]);
	let layout = familyTreeDAG.getLayoutedElements(
		data.Nodes,
		data.Edges,
		tailwindClassToPixels('w-40') || 160,
		tailwindClassToPixels('h-40') || 160,
		'TB'
	);
	nodes = layout.Nodes;
	edges = layout.Edges;

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
				if (Number(data.id) === Number(node.id)) {
					relationshipStart = null;
					openPersonMenu = undefined;

					return;
				}
				fetch('/api/person/' + node.id, {
					method: 'DELETE',
					headers: {
						'Content-Type': 'application/json'
					}
				})
					.then((response) => {
						if (response.ok) {
							nodes = nodes.filter((n) => n.id !== node.id);
							edges = edges.filter((e) => e.source !== node.id && e.target !== node.id);
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
				relationshipStart = Number(node.id);
				createPerson = true;
				openPersonMenu = undefined;
			},
			addRelationship: () => {
				relationshipStart = Number(node.id);
				openPersonMenu = undefined;
			},
			addAdmin: () => {
				relationshipStart = Number(node.id);
				openPersonMenu = undefined;
			},
			addRecipe: () => {
				relationshipStart = Number(node.id);
				openPersonMenu = undefined;
			},
			id: node.id,
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

		let newLayout = familyTreeDAG.getLayoutedElements(
			nodes,
			edges,
			tailwindClassToPixels('w-40') || 160,
			tailwindClassToPixels('h-40') || 160,
			'TB'
		);

		edges = newLayout.Edges;
		nodes = newLayout.Nodes;
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
		const sourceNodeId = connectionState.fromNode?.id;
		if (sourceNodeId === undefined) return;
		relationshipStart = Number(sourceNodeId);
		createPerson = true;
		console.log('createPerson', createPerson);
		console.log('relationshipStart', relationshipStart);
	};
</script>

<svelte:head>
	<title>{title({ page: family_tree() })}</title>
</svelte:head>
<div class="drawer">
	<div style="height:100vh;" class="!bg-base-200 drawer-content flex flex-col">
		<SvelteFlowProvider>
			<SvelteFlow
				bind:nodes
				bind:edges
				onconnectend={handleConnectEnd}
				onnodeclick={handleNodeClickFunc}
				onnodecontextmenu={handleContextMenu}
				onpaneclick={handlePaneClick}
				class="!bg-base-200"
				{nodeTypes}
				{edgeTypes}
				fitView
				onlyRenderVisibleElements
				connectionLineType={ConnectionLineType.SmoothStep}
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
	<SideBar />
</div>
