<script lang="ts">
	import { onMount } from 'svelte';
	import {
		variations as variationsLabel,
		no_variations,
		new_variation,
		variation_notes,
		variation_by,
		recipe as recipeLabel,
		save,
		cancel
	} from '$lib/paraglide/messages';
	import { languageTag } from '$lib/paraglide/runtime';
	import { createVariation, getVariations } from './api';
	import { errorMessage } from './errors';
	import { cleanDraft, formatUnixSeconds, personName, toDraft } from './model';
	import type { RecipeProperties, RecipeVariation } from './model';
	import RecipeFields from './RecipeFields.svelte';

	let {
		recipeId,
		baseProps,
		onOpen
	}: {
		recipeId: number;
		baseProps: RecipeProperties;
		onOpen: (id: number) => void;
	} = $props();

	let variations: RecipeVariation[] = $state([]);
	let isLoading = $state(true);
	let busy = $state(false);
	let error: string | null = $state(null);
	let showForm = $state(false);
	let draft = $state(toDraft(null));
	let notes = $state('');

	async function load() {
		const result = await getVariations(recipeId);
		if (result.ok) {
			variations = result.data.variations ?? [];
		} else {
			error = errorMessage(result.status);
		}
		isLoading = false;
	}

	onMount(load);

	function openForm() {
		draft = toDraft(baseProps);
		notes = '';
		error = null;
		showForm = true;
	}

	async function submit() {
		if (busy || draft.name.trim() === '') return;
		busy = true;
		error = null;
		const result = await createVariation(recipeId, cleanDraft(draft), notes);
		if (result.ok) {
			showForm = false;
			await load();
		} else {
			error = errorMessage(result.status);
		}
		busy = false;
	}
</script>

<section class="flex flex-col gap-2">
	<div class="flex items-center justify-between">
		<h4 class="text-base font-bold">{variationsLabel()}</h4>
		{#if !showForm}
			<button class="btn btn-accent btn-xs" onclick={openForm}>{new_variation()}</button>
		{/if}
	</div>

	{#if error}
		<div class="alert alert-error alert-soft py-2 text-sm" role="alert">{error}</div>
	{/if}

	{#if showForm}
		<div class="bg-base-200 rounded-box flex flex-col gap-3 p-3">
			<div>
				<span class="label font-semibold">{variation_notes()}</span>
				<textarea class="textarea textarea-bordered w-full" rows="2" bind:value={notes}></textarea>
			</div>
			<RecipeFields bind:draft />
			<div class="flex justify-end gap-2">
				<button class="btn btn-ghost btn-sm" onclick={() => (showForm = false)}>
					{cancel()}
				</button>
				<button
					class="btn btn-primary btn-sm"
					onclick={submit}
					disabled={busy || draft.name.trim() === ''}
				>
					{save()}
				</button>
			</div>
		</div>
	{:else if isLoading}
		<span class="loading loading-spinner"></span>
	{:else if variations.length === 0}
		<p class="text-base-content/60 text-sm">{no_variations()}</p>
	{:else}
		<ul class="flex flex-col gap-2">
			{#each variations as v (v.variation?.Id)}
				<li>
					<button
						class="bg-base-200 hover:bg-base-300 rounded-box w-full p-3 text-left transition-colors"
						onclick={() => v.variation?.Id !== undefined && onOpen(v.variation.Id)}
					>
						<div class="font-semibold">{v.variation?.Props?.name ?? recipeLabel()}</div>
						<div class="text-base-content/60 text-xs">
							{variation_by({ name: personName(v.creator) })}
							{#if v.variation_relationship?.created_at}
								- {formatUnixSeconds(v.variation_relationship.created_at, languageTag())}
							{/if}
						</div>
						{#if v.variation_relationship?.notes}
							<p class="mt-1 line-clamp-2 text-sm">{v.variation_relationship.notes}</p>
						{/if}
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</section>
