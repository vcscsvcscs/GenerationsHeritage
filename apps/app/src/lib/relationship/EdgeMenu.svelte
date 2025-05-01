<script lang="ts">
	import { onMount } from 'svelte';
	import {
		add_relationship,
		remove,
		create_relationship_and_person,
		add_administrator
	} from '$lib/paraglide/messages';
	import type { Edge } from '@xyflow/svelte';

	export let edge: Edge;
	export let XUserId: string;
	export let top: number | undefined;
	export let left: number | undefined;
	export let right: number | undefined;
	export let bottom: number | undefined;
	export let onClick: () => void;
	export let deleteEdge: () => void;

	let contextMenu: HTMLDivElement;
	let isAdmin: boolean = false;
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

		fetch(`/api/admin/${edge.source}/${XUserId}`)
			.then((response) => {
				if (response.status === 200) {
					isAdmin = true;
				} else {
					isAdmin = false;
				}
			})
			.catch((error) => {
				console.error('Error fetching admin status:', error);
			});
		fetch(`/api/admin/${edge.target}/${XUserId}`)
			.then((response) => {
				if (response.status === 200) {
					isAdmin = true;
				}
			})
			.catch((error) => {
				console.error('Error fetching admin status:', error);
			});
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
	{#if isAdmin}
		<button onclick={deleteEdge} class="btn">{remove()}</button>
	{/if}
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
