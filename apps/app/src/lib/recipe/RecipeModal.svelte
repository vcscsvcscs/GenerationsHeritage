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
		remove,
		favourite
	} from '$lib/paraglide/messages';
	import type { components } from '$lib/api/api.gen';

	let {
		recipeData,
		recipeId,
		editable = false,
		liked = false,
		closeModal,
		onSaved,
		onLikeToggle
	}: {
		recipeData: components['schemas']['RecipeProperties'];
		recipeId: number | undefined;
		editable?: boolean;
		liked?: boolean;
		closeModal: () => void;
		onSaved?: () => void;
		onLikeToggle?: (liked: boolean) => void;
	} = $props();

	let editorMode = $state(false);
	let isLiked = $state(liked);
	let likeLoading = $state(false);
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

	async function toggleLike() {
		if (!recipeId || likeLoading) return;
		likeLoading = true;
		try {
			if (isLiked) {
				const response = await fetch(`/api/recipe/${recipeId}/relationship?personId=me`, {
					method: 'DELETE'
				});
				if (response.ok) {
					isLiked = false;
					onLikeToggle?.(false);
				}
			} else {
				const response = await fetch(`/api/recipe/${recipeId}/relationship`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						id: 0,
						relationship: { schema: { like_it: true } }
					})
				});
				if (response.ok) {
					isLiked = true;
					onLikeToggle?.(true);
				}
			}
		} catch (e) {
			console.error('Error toggling like:', e);
		} finally {
			likeLoading = false;
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
				<div class="flex items-center gap-2">
					<h3 class="text-lg font-bold">{recipeData.name ?? recipe()}</h3>
					{#if onLikeToggle}
						<button
							class="btn btn-ghost btn-sm"
							class:text-error={isLiked}
							onclick={toggleLike}
							disabled={likeLoading}
							title={favourite()}
						>
							{#if isLiked}
								<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5">
									<path d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3A5.5 5.5 0 0112 5.052 5.5 5.5 0 0116.313 3c2.973 0 5.437 2.322 5.437 5.25 0 3.925-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001a.752.752 0 01-.704 0l-.003-.001z" />
								</svg>
							{:else}
								<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
									<path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
								</svg>
							{/if}
						</button>
					{/if}
				</div>
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
