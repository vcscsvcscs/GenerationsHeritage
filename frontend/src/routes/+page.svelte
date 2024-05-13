<script lang="ts">
	import { writable } from 'svelte/store';
	import { onMount } from 'svelte';
	import { SvelteFlowProvider, SvelteFlow, Controls, MiniMap, Background } from '@xyflow/svelte';
	import type { Node, Edge, NodeTypes } from '@xyflow/svelte';
	import { isAuthenticated } from '../lib/stores';
	import PersonNode from './../lib/family_tree/PersonNode.svelte';
	import { setFamilyTreeNodes } from '../lib/family_tree/setFamilyTreeNodes';
	import { login } from '../lib/auth';

	import PersonMenu from '$lib/family_tree/PersonMenu.svelte';
	import PersonPanel from '$lib/family_tree/PersonPanel.svelte';
	const nodes = writable<Node[]>([]);
	const edges = writable<Edge[]>([]);
	onMount(() => {
		if (!$isAuthenticated) {
			console.log('user is not authenticated');
			login();
		} else {
			console.log('user is authenticated');
			if (!setFamilyTreeNodes()){
				
			};
		}
	});

	const nodeTypes: NodeTypes = {
		custom: PersonNode
	};

	let personPanel: {} | null;

	let menu: { id: string; top?: number; left?: number; right?: number; bottom?: number } | null;
	let width: number;
	let height: number;

	function handleContextMenu({ detail: { event, node } }) {
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
			{#if menu}
				<PersonMenu
					onClick={handlePaneClick}
					id={menu.id}
					top={menu.top}
					left={menu.left}
					right={menu.right}
					bottom={menu.bottom}
				/>
			{/if}
			<PersonPanel showPanel={personPanel != null} data={personPanel} />
		</SvelteFlow>
	</SvelteFlowProvider>
</div>
