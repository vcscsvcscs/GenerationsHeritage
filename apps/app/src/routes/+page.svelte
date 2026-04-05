<script lang="ts">
	import CreateRelationship from '$lib/relationship/Modal.svelte';
	import { onMount } from 'svelte';
	import { nodeTypes, edgeTypes } from '$lib/graph/model';
	import { title, family_tree, select } from '$lib/paraglide/messages.js';
	import type { RelationshipMenu } from '$lib/relationship/model.ts';
	import AdminMenu from '$lib/admin/Modal.svelte';

	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import type { OnConnectEnd, Node, Edge, NodeEventWithPointer } from '@xyflow/svelte';

	import PersonModal from '$lib/profile/Modal.svelte';
	import PersonMenu from '$lib/graph/PersonMenu.svelte';
	import CreatePerson from '$lib/profile/create/Modal.svelte';

	import type { components } from '$lib/api/api.gen';
	import type { NodeMenu } from '$lib/graph/model';
	import RecipeListModal from '$lib/recipe/RecipeListModal.svelte';
	import CookbookModal from '$lib/recipe/CookbookModal.svelte';
	import { cookbook } from '$lib/paraglide/messages.js';

	import { handleNodeClick } from '$lib/graph/node_click';

	import { FamilyTree } from '$lib/graph/layout';
	import { tailwindClassToPixels } from '$lib/tailwindSizeToPx';
	import type { Layout } from '$lib/graph/model';
	import HamburgerIcon from '$lib/sidebar/hamburgerIcon.svelte';

	let { data }: { data: Layout & { id: string } } = $props();

	let selectedPerson: components['schemas']['PersonProperties'] & { id: string | undefined } =
		$state({
			id: undefined
		});
	let selectedRelationship: Edge | undefined = $state(undefined);
	let openPersonPanel = $state(false);
	let openPersonMenu: NodeMenu | undefined = $state(undefined);
	let with_out_spouse = $state(false);
	let createRelationship = $state(false);
	let adminMenu = $state(false);
	let recipePersonId: number | null = $state(null);
	let recipePersonName = $state('');
	let showCookbook = $state(false);

	let familyTreeDAG = new FamilyTree();
	let layout = familyTreeDAG.getLayoutedElements(
		data.Nodes,
		data.Edges,
		tailwindClassToPixels('w-40') || 160,
		tailwindClassToPixels('h-40') || 160,
		'TB'
	);
	let nodes = $state.raw<Node[]>([] as Node[]);
	let edges = $state.raw<Edge[]>([] as Edge[]);

	let relationshipStart: number | null = $state(null);
	let relationshipMenu = $state(undefined as RelationshipMenu | undefined);
	let createPerson = $state(false);

	let clientWidth: number | undefined = $state();
	let clientHeight: number | undefined = $state();

	let removePersonFromGraph = (id: any) => {
		nodes = nodes.filter((n) => n.data.id !== id);
		edges = edges.filter((e) => e.source !== 'person' + id && e.target !== 'person' + id);
	};

	let delete_profile = (id: any) => {
		fetch('/api/person/' + id, {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json'
			}
		})
			.then((response) => {
				if (response.ok) {
					removePersonFromGraph(id);
				} else {
					alert('Error deleting person');
				}
			})
			.catch((error) => {
				console.error('Error:', error);
			});
	};

	const handleContextMenu: NodeEventWithPointer<MouseEvent> = ({ event, node }) => {
		event.preventDefault();

		if (clientHeight === undefined || clientWidth === undefined) {
			clientHeight = window.innerHeight;
			clientWidth = window.innerWidth;
		}

		if (openPersonMenu !== undefined) {
			openPersonMenu.onClick();
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

				delete_profile(node.data.id);
				openPersonMenu = undefined;
			},
			createRelationshipAndNode: () => {
				relationshipStart = Number(node.data.id);
				createPerson = true;
				openPersonMenu = undefined;
			},
			addRelationship: () => {
				relationshipStart = Number(node.data.id);
				createRelationship = true;
				selectedRelationship = {
					id: 'relationship' + node.data.id,
					source: String(relationshipStart),
					target: String(node.data.id)
				};
				openPersonMenu = undefined;
			},
			addAdmin: () => {
				relationshipStart = Number(node.data.id);
				openPersonMenu = undefined;
			},
			addRecipe: () => {
				recipePersonId = Number(node.data.id);
				recipePersonName = [node.data.first_name, node.data.last_name].filter(Boolean).join(' ');
				openPersonMenu = undefined;
			},
			id: String(node.data.id),
			top: event.clientY < clientHeight - 200 ? event.clientY : undefined,
			left: event.clientX < clientWidth - 200 ? event.clientX : undefined,
			right: event.clientX >= clientWidth - 200 ? clientWidth - event.clientX : undefined,
			bottom: event.clientY >= clientHeight - 200 ? clientHeight - event.clientY : undefined
		};
	};

	function onCreation(newNodes: Array<Node> | null, newEdges: Array<Edge> | null): void {
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
		edges = [...newLayout.Edges];
		nodes = [...newLayout.Nodes];
	}

	let handleNodeClickFunc = handleNodeClick(
		(
			person: components['schemas']['PersonProperties'] & {
				id: number | undefined;
			}
		) => {
			selectedPerson = { ...person, id: String(person.id) };
			openPersonPanel = true;
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
							id: string | undefined;
						};
						selectedPerson.id = String(person.id);
					}
					console.debug('Fetched person data:', data);
				});
		}
	);

	let handlePaneClick = ({ event }: { event: MouseEvent }) => {
		openPersonPanel = false;
		openPersonMenu = undefined;
	};

	const handleConnectEnd: OnConnectEnd = (event, connectionState) => {
		event.preventDefault();
		const sourceNodeId = connectionState.fromNode?.data.id;
		if (sourceNodeId === undefined) return;
		relationshipStart = Number(sourceNodeId);
		if (connectionState.isValid) {
			createRelationship = true;
			selectedRelationship = {
				id: 'relationship' + connectionState.toNode?.data.id,
				source: String(relationshipStart),
				target: String(connectionState.toNode?.data.id)
			};
			return;
		}

		createPerson = true;
	};
	onMount(() => {
		nodes = [...layout.Nodes];
		edges = [...layout.Edges];
	});
</script>

<svelte:head>
	<title>{title({ page: family_tree() })}</title>
</svelte:head>
<div style="height:100vh;" class="!bg-base-200 flex flex-col">
	<SvelteFlowProvider>
		<SvelteFlow
			bind:nodes
			bind:edges
			onconnectend={handleConnectEnd}
			onedgeclick={({ edge, event }: { edge: Edge; event: MouseEvent }) => {
				selectedRelationship = edge;
				selectedRelationship.source = String(edge.source.replace('person', ''));
				selectedRelationship.target = String(edge.target.replace('person', ''));
			}}
			onnodeclick={handleNodeClickFunc}
			onnodecontextmenu={handleContextMenu}
			onedgecontextmenu={({ edge, event }: { edge: Edge; event: MouseEvent }) => {
				selectedRelationship = edge;
				selectedRelationship.source = String(edge.source.replace('person', ''));
				selectedRelationship.target = String(edge.target.replace('person', ''));
				if (clientHeight === undefined || clientWidth === undefined) {
					clientHeight = window.innerHeight;
					clientWidth = window.innerWidth;
				}
				relationshipMenu = {
					XUserId: data.id,
					edge: selectedRelationship,
					onClick: () => {
						relationshipMenu = undefined;
					},
					deleteEdge: () => {
						edges = edges.filter((e) => e.id !== edge.id);
						relationshipMenu = undefined;
					},
					top: event.clientY < clientHeight - 200 ? event.clientY : undefined,
					left: event.clientX < clientWidth - 200 ? event.clientX : undefined,
					right: event.clientX >= clientWidth - 200 ? clientWidth - event.clientX : undefined,
					bottom: event.clientY >= clientHeight - 200 ? clientHeight - event.clientY : undefined
				};
			}}
			onpaneclick={handlePaneClick}
			class="!bg-base-200"
			{nodeTypes}
			{edgeTypes}
			fitView={true}
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
					onCreation={(node, edges) => {
						onCreation([node], edges);
						createPerson = false;
					}}
					closeModal={() => {
						createPerson = false;
					}}
					relationshipStartID={relationshipStart}
				></CreatePerson>
			{/if}
			{#if selectedRelationship}
				<CreateRelationship
					{createRelationship}
					onCreation={(newEdges: Array<Edge>) => {
						onCreation(null, newEdges);
						createRelationship = false;
					}}
					closeModal={() => {
						createRelationship = false;
						selectedRelationship = undefined;
						relationshipStart = null;
						layout = familyTreeDAG.getLayoutedElements(
							nodes,
							edges,
							tailwindClassToPixels('w-40') || 160,
							tailwindClassToPixels('h-40') || 160,
							'TB'
						);
						edges = [...layout.Edges];
						nodes = [...layout.Nodes];
					}}
					startNode={String(selectedRelationship.source)}
					endNode={String(selectedRelationship.target)}
				/>
			{/if}
			{#if openPersonMenu !== undefined}
				<PersonMenu {...openPersonMenu!} />
			{/if}
			{#if recipePersonId !== null}
				<RecipeListModal
					personId={recipePersonId}
					personName={recipePersonName}
					closeModal={() => {
						recipePersonId = null;
						recipePersonName = '';
					}}
				/>
			{/if}
			{#if showCookbook}
				<CookbookModal
					closeModal={() => {
						showCookbook = false;
					}}
				/>
			{/if}
			{#if adminMenu}
				<AdminMenu
					createProfile={() => {
						createPerson = true;
						relationshipStart = null;
					}}
					createRelationshipAndProfile={(id: number) => {
						createPerson = true;
						relationshipStart = id;
					}}
					addRelationship={(id: number) => {
						createRelationship = true;
						selectedRelationship = {
							id: 'relationship' + id,
							source: String(id),
							target: String(id)
						};
					}}
					closeModal={() => {
						adminMenu = false;
					}}
					editProfile={(id: number) => {
						selectedPerson = { id: String(id) };
						fetch('/api/person/' + id, {
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
										id: string | undefined;
									};
									selectedPerson.id = String(id);
									openPersonPanel = true;
								} else {
									alert('Error fetching person data');
								}
							});
					}}
					removePersonFromGraph={removePersonFromGraph}
				/>
			{/if}
		</SvelteFlow>
	</SvelteFlowProvider>
</div>

<div class="absolute left-2 top-2 flex flex-row items-center gap-2">
	<HamburgerIcon
		open_admin_panel={() => {
			adminMenu = !adminMenu;
		}}
	/>
	<button
		class="btn btn-sm btn-primary"
		onclick={() => {
			showCookbook = !showCookbook;
		}}
	>
		{cookbook()}
	</button>
</div>
