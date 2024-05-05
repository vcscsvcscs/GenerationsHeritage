<script lang="ts">
	export let onClick: () => void;
	export let deleteNode: () => void;
	export let addRelationship: () => void;
	export let addFamilymember: () => void;
	export let id: string;
	export let top: number | undefined;
	export let left: number | undefined;
	export let right: number | undefined;
	export let bottom: number | undefined;
</script>

<div
	style="top: {top}px; left: {left}px; right: {right}px; bottom: {bottom}px;"
	class="context-menu bg-primary-100 rounded-lg shadow-lg"
	on:click={onClick}
>
	<p style="margin: 0.5em;">
		<small>node: {id}</small>
	</p>
	<button on:click={addRelationship}>Add Relationship</button>
	<button on:click={addFamilymember}>Add familymember</button>
	<button on:click={deleteNode}>Delete</button>
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
