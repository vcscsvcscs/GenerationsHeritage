<script lang="ts">
	import { back, biography, close, edit } from '$lib/paraglide/messages';

	export let editorMode = false;

	export let onClose: () => void;
	export let onToggleEdit: () => void;
</script>

<div class="flex items-center justify-between p-2">
	<h3 class="text-lg font-bold">{biography()}</h3>
	<div class="space-x-2">
		<button class="btn btn-outline btn-error btn-sm" on:click={onClose}>
			{close()}
		</button>
		<button class="btn btn-outline btn-primary btn-sm" on:click={onToggleEdit}>
			{editorMode ? back() : edit}
		</button>
	</div>
</div>
