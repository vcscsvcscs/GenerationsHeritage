<script lang="ts">
	import { onMount } from 'svelte';
	import {
		add_relationship,
		remove,
		create_relationship_and_person,
		add_administrator
	} from '$lib/paraglide/messages';

	export let top: number | undefined;
	export let left: number | undefined;
	export let right: number | undefined;
	export let bottom: number | undefined;
	export let onClick: () => void;
	export let deleteNode: () => void;
	export let createRelationshipAndNode: () => void;
	export let addRelationship: () => void;
	export let addAdmin: (() => void) | undefined;

	let contextMenu: HTMLDivElement;
	onMount(() => {
		if (top) {
			contextMenu.style.top = `${top}px`;
		}
		if (left) {
			contextMenu.style.left = `${left}px`;
		}
		if (right) {
			contextMenu.style.right = `${right}px`;
		}
		if (bottom) {
			contextMenu.style.bottom = `${bottom}px`;
		}
	});
</script>

<div
	role="menu"
	tabindex="-1"
	bind:this={contextMenu}
	class="context-menu bg-primary-100 rounded-lg shadow-lg"
	onclick={onClick}
	onkeydown={(e) => {
		if (e.key === 'Esc' || e.key === ' ' || e.key === 'Escape') {
			onClick();
		}
	}}
>
	<button onclick={createRelationshipAndNode} class="btn">
		{create_relationship_and_person()}
	</button>
	<button onclick={addRelationship} class="btn">{add_relationship()}</button>
	<button onclick={addAdmin} class="btn">{add_administrator()}</button>
	<button onclick={deleteNode} class="btn">{remove()}</button>
</div>

<style>
	.context-menu {
		border-style: solid;
		box-shadow: 10px 19px 20px rgba(0, 0, 0, 10%);
		position: absolute;
		z-index: 10;
	}

	.context-menu button {
		border: none;
		display: block;
		padding: 0.5em;
		text-align: left;
		width: 100%;
	}
</style>
