<script lang="ts">
	import { writable } from 'svelte/store';
	import { onMount } from 'svelte';
	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';
	import type { Node, Edge, NodeTypes } from '@xyflow/svelte';
	import { isAuthenticated } from '../lib/stores';
	import PersonNode from './../lib/family_tree/PersonNode.svelte';
	import { setFamilyTreeNodes } from '../lib/family_tree/setFamilyTreeNodes';
	import { login } from '../lib/auth';

	import PersonMenu from '$lib/family_tree/PersonMenu.svelte';
	import PersonPanel from '$lib/family_tree/PersonPanel.svelte';
	import AddRelationship from '$lib/family_tree/AddRelationship.svelte';
	import AddFamilyMember from '$lib/family_tree/AddFamilyMember.svelte';
	import CreateProfile from '$lib/family_tree/CreateProfile.svelte';

	let relationshipPanel = '';
	let AddFamilyMemberPanel = '';
	let CreateProfilePanel = false;

	const nodes = writable<Node[]>([]);
	const edges = writable<Edge[]>([]);
	onMount(() => {
		if (!$isAuthenticated) {
			console.log('user is not authenticated');
			login();
		} else {
			console.log('user is authenticated');
			if (!setFamilyTreeNodes()) {
				CreateProfilePanel = true;
			}
		}
	});

	const nodeTypes: NodeTypes = {
		custom: PersonNode
	};

	let personPanel:
		| {
				id: string;
				last_name: string;
				first_name: string;
				middle_name: string;
				profile_picture: string;
		  }
		| undefined;

	let menu: { id: string; top?: number; left?: number; right?: number; bottom?: number } | null;
	let width: number;
	let height: number;

	function handleContextMenu({
		detail: { event, node }
	}: {
		detail: { event: MouseEvent; node: Node };
	}) {
		event.preventDefault();

		menu = {
			id: node.data.ID,
			top: event.clientY < height - 200 ? event.clientY : undefined,
			left: event.clientX < width - 200 ? event.clientX : undefined,
			right: event.clientX >= width - 200 ? width - event.clientX : undefined,
			bottom: event.clientY >= height - 200 ? height - event.clientY : undefined
		};
	}

	// Close the context menu if it's open whenever the window is clicked.
	function handlePaneClick() {
		menu = null;
		relationshipPanel = '';
	}
</script>

<div style="height:100vh;">
	<SvelteFlowProvider>
		<SvelteFlow
			{nodes}
			{nodeTypes}
			{edges}
			on:nodecontextmenu={handleContextMenu}
			on:nodeclick
			on:paneclick={handlePaneClick}
			class="bg-base-100"
			fitView
			onlyRenderVisibleElements
		>
			<Controls class="bg-base-300 text-base-content" />
			<MiniMap class="bg-base-200" />
			{#if menu != null}
				<PersonMenu
					onClick={handlePaneClick}
					deleteNode={() => {
						console.log('delete node');
					}}
					addRelationship={() => {
						if (menu) relationshipPanel = menu.id;
					}}
					addFamilymember={() => {
						if (menu) AddFamilyMemberPanel = menu.id;
					}}
					id={menu.id}
					top={menu.top}
					left={menu.left}
					right={menu.right}
					bottom={menu.bottom}
				/>
			{/if}
			{#if personPanel != null}
				<PersonPanel data={personPanel} />
			{/if}
			{#if relationshipPanel != ''}
				<AddRelationship id={relationshipPanel} />
			{/if}
			{#if AddFamilyMemberPanel != ''}
				<AddFamilyMember id={AddFamilyMemberPanel} />
			{/if}
			{#if CreateProfilePanel}
				<CreateProfile />
			{/if}
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
