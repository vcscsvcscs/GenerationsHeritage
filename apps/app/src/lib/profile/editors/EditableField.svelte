<script lang="ts">
	import type { components } from '$lib/api/api.gen';

	export let key: keyof components['schemas']['PersonProperties'];
	export let value: any;
	export let editorMode = false;
	export let onChange: (field: keyof components['schemas']['PersonProperties'], value: any) => void;
	let numberField: HTMLInputElement;
	let textField: HTMLTextAreaElement;
	let checkboxField: HTMLInputElement;
</script>

{#if editorMode}
	{#if typeof value === 'boolean'}
		<input
			type="checkbox"
			class="toggle toggle-primary"
			checked={value}
			bind:this={checkboxField}
			oninput={() => onChange(key, checkboxField.value === 'true')}
		/>
	{:else if typeof value === 'number'}
		<input
			type="number"
			class="input input-bordered input-sm w-full"
			{value}
			bind:this={numberField}
			oninput={() => onChange(key, Number(numberField.value))}
		/>
	{:else}
		<textarea
			class="textarea textarea-bordered textarea-sm w-full"
			{value}
			oninput={() => onChange(key, textField.value)}
			bind:this={textField}
		></textarea>
	{/if}
{:else}
	<p class="text-sm text-gray-700">{value ?? '-'}</p>
{/if}
