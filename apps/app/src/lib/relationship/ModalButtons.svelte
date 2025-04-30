<script lang="ts">
	import { add_relationship, back, biography, close, create, edit, save } from '$lib/paraglide/messages';

	export let editorMode = false;
	export let createMode = false;

	export let onClose: () => void;
	export let onToggleEdit: () => void;
	export let onSave: () => void;
	export let onCreate: () => void;
</script>

<div class="flex items-center justify-between p-2">
	<h3 class="text-lg font-bold">{biography()}</h3>
	<div class="space-x-2">
		{#if !createMode}
			<button class="btn btn-secondary btn-sm" on:click={onToggleEdit}>
				{editorMode ? back() : edit()}
			</button>
		{/if}
		{#if editorMode}
			{#if createMode}
				<button class="btn btn-accent btn-sm" on:click={onCreate}>
					{add_relationship()}
				</button>
			{:else}
				<button class="btn btn-accent btn-sm" on:click={onSave}>
					{save()}
				</button>
			{/if}
		{/if}
		<button class="btn btn-error btn-sm" on:click={onClose}>
			{close()}
		</button>
	</div>
</div>
