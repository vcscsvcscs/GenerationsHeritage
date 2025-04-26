<script lang="ts">
	import { onMount } from 'svelte';
	import { type NodeMenu } from '$lib/graph/model';
	import {
		add_relationship,
		remove,
		create_relationship_and_person,

		add_administrator

	} from '$lib/paraglide/messages';

	let props: NodeMenu = $props();

	let contextMenu: HTMLDivElement;
	onMount(() => {
		if (props.top) {
			contextMenu.style.top = `${props.top}px`;
		}
		if (props.left) {
			contextMenu.style.left = `${props.left}px`;
		}
		if (props.right) {
			contextMenu.style.right = `${props.right}px`;
		}
		if (props.bottom) {
			contextMenu.style.bottom = `${props.bottom}px`;
		}
	});
</script>

<div
	role="menu"
	tabindex="-1"
	bind:this={contextMenu}
	class="context-menu bg-primary-100 rounded-lg shadow-lg"
	onclick={props.onClick}
	onkeydown={(e) => {
		if (e.key === 'Esc' || e.key === ' ' || e.key === 'Escape') {
			props.onClick();
		}
	}}
>
	<button onclick={props.createRelationshipAndNode} class="btn"
		>{create_relationship_and_person()}</button
	>
	<button onclick={props.addRelationship} class="btn">{add_relationship()}</button>
	<button onclick={props.addAdmin} class="btn">{add_administrator()}</button>
	<button onclick={props.deleteNode} class="btn">{remove()}</button>
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
