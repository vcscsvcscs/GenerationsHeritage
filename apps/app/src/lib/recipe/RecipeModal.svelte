<script lang="ts">
	import { fade } from 'svelte/transition';
	import {
		close,
		edit,
		back,
		save,
		recipe,
		description,
		origin,
		category,
		ingredients,
		instructions,
		notes,
		add,
		remove
	} from '$lib/paraglide/messages';
	import type { components } from '$lib/api/api.gen';

	let {
		recipeData,
		recipeId,
		editable = false,
		closeModal,
		onSaved
	}: {
		recipeData: components['schemas']['RecipeProperties'];
		recipeId: number | undefined;
		editable?: boolean;
		closeModal: () => void;
		onSaved?: () => void;
	} = $props();

	let editorMode = $state(false);
	let draft = $state({
		...recipeData,
		ingredients: recipeData.ingredients ?? [],
		instructions: recipeData.instructions ?? []
	});

	function toggleEdit() {
		if (!editorMode) {
			draft = {
				...recipeData,
				ingredients: recipeData.ingredients ?? [],
				instructions: recipeData.instructions ?? []
			};
		}
		editorMode = !editorMode;
	}

	async function handleSave() {
		if (!recipeId) return;
		try {
			const response = await fetch(`/api/recipe/${recipeId}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(draft)
			});
			if (response.ok) {
				Object.assign(recipeData, draft);
				editorMode = false;
				onSaved?.();
			} else {
				alert('Error saving recipe');
			}
		} catch (e) {
			alert('Error saving recipe: ' + e);
		}
	}

	function addIngredient() {
		draft.ingredients = [...draft.ingredients, ''];
	}

	function removeIngredient(index: number) {
		draft.ingredients = draft.ingredients.filter((_, i) => i !== index);
	}

	function addInstruction() {
		draft.instructions = [...draft.instructions, ''];
	}

	function removeInstruction(index: number) {
		draft.instructions = draft.instructions.filter((_, i) => i !== index);
	}
</script>

<div class="modal modal-open" transition:fade>
	<div class="modal-box max-h-screen w-full max-w-3xl overflow-y-auto">
		<div class="bg-base-100 z-7 sticky top-0">
			<div class="flex items-center justify-between p-2">
				<h3 class="text-lg font-bold">{recipeData.name ?? recipe()}</h3>
				<div class="space-x-2">
					{#if editable}
						<button class="btn btn-secondary btn-sm" onclick={toggleEdit}>
							{editorMode ? back() : edit()}
						</button>
						{#if editorMode}
							<button class="btn btn-accent btn-sm" onclick={handleSave}>
								{save()}
							</button>
						{/if}
					{/if}
					<button class="btn btn-error btn-sm" onclick={closeModal}>
						{close()}
					</button>
				</div>
			</div>
			<div class="divider"></div>
		</div>

		<div class="flex flex-col gap-4 p-2">
			<!-- Name -->
			<div>
				<label class="label font-semibold">{recipe()}</label>
				{#if editorMode}
					<input type="text" class="input input-bordered w-full" bind:value={draft.name} />
				{:else}
					<p class="text-lg">{recipeData.name ?? '-'}</p>
				{/if}
			</div>

			<!-- Origin & Category -->
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div>
					<label class="label font-semibold">{origin()}</label>
					{#if editorMode}
						<input type="text" class="input input-bordered w-full" bind:value={draft.origin} />
					{:else}
						<p>{recipeData.origin ?? '-'}</p>
					{/if}
				</div>
				<div>
					<label class="label font-semibold">{category()}</label>
					{#if editorMode}
						<input type="text" class="input input-bordered w-full" bind:value={draft.category} />
					{:else}
						<p>{recipeData.category ?? '-'}</p>
					{/if}
				</div>
			</div>

			<!-- Description -->
			<div>
				<label class="label font-semibold">{description()}</label>
				{#if editorMode}
					<textarea class="textarea textarea-bordered w-full" rows="3" bind:value={draft.description}></textarea>
				{:else}
					<p>{recipeData.description ?? '-'}</p>
				{/if}
			</div>

			<!-- Ingredients -->
			<div>
				<label class="label font-semibold">{ingredients()}</label>
				{#if editorMode}
					{#each draft.ingredients ?? [] as item, i}
						<div class="flex items-center gap-2 mb-1">
							<input
								type="text"
								class="input input-bordered input-sm flex-1"
								bind:value={draft.ingredients[i]}
							/>
							<button class="btn btn-xs btn-ghost text-error" onclick={() => removeIngredient(i)}>
								&#10005;
							</button>
						</div>
					{/each}
					<button class="btn btn-accent btn-xs mt-1" onclick={addIngredient}>
						{add()}
					</button>
				{:else}
					{#if recipeData.ingredients && recipeData.ingredients.length > 0}
						<ul class="list-disc pl-5">
							{#each recipeData.ingredients as item}
								<li>{item}</li>
							{/each}
						</ul>
					{:else}
						<p>-</p>
					{/if}
				{/if}
			</div>

			<!-- Instructions -->
			<div>
				<label class="label font-semibold">{instructions()}</label>
				{#if editorMode}
					{#each draft.instructions ?? [] as step, i}
						<div class="flex items-center gap-2 mb-1">
							<span class="text-sm font-mono w-6">{i + 1}.</span>
							<textarea
								class="textarea textarea-bordered textarea-sm flex-1"
								bind:value={draft.instructions[i]}
							></textarea>
							<button class="btn btn-xs btn-ghost text-error" onclick={() => removeInstruction(i)}>
								&#10005;
							</button>
						</div>
					{/each}
					<button class="btn btn-accent btn-xs mt-1" onclick={addInstruction}>
						{add()}
					</button>
				{:else}
					{#if recipeData.instructions && recipeData.instructions.length > 0}
						<ol class="list-decimal pl-5">
							{#each recipeData.instructions as step}
								<li class="mb-1">{step}</li>
							{/each}
						</ol>
					{:else}
						<p>-</p>
					{/if}
				{/if}
			</div>

			<!-- Notes -->
			<div>
				<label class="label font-semibold">{notes()}</label>
				{#if editorMode}
					<textarea class="textarea textarea-bordered w-full" rows="2" bind:value={draft.notes}></textarea>
				{:else}
					<p>{recipeData.notes ?? '-'}</p>
				{/if}
			</div>
		</div>
	</div>
</div>
