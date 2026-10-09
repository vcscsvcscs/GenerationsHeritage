<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import {
		close,
		edit,
		back,
		save,
		cancel,
		recipe,
		description,
		origin,
		category,
		ingredients,
		instructions,
		notes,
		favourite,
		delete_recipe,
		delete_recipe_confirm,
		delete_action,
		loading
	} from '$lib/paraglide/messages';
	import { deleteRecipe, getRecipe, likeRecipe, unlikeRecipe, updateRecipe } from './api';
	import { errorMessage } from './errors';
	import { canEdit as entryCanEdit, cleanDraft, toDraft } from './model';
	import type { RecipeProperties } from './model';
	import RecipeFields from './RecipeFields.svelte';
	import RecipeComments from './RecipeComments.svelte';
	import RecipeVariations from './RecipeVariations.svelte';

	type View = { id: number; props: RecipeProperties; canEdit: boolean };

	let {
		recipeId,
		initial,
		currentUserId = null,
		liked = false,
		closeModal,
		onSaved,
		onLikeToggle
	}: {
		recipeId: number;
		initial?: { recipe: RecipeProperties; canEdit: boolean };
		currentUserId?: number | null;
		liked?: boolean;
		closeModal: () => void;
		onSaved?: () => void;
		onLikeToggle?: (liked: boolean) => void;
	} = $props();

	let trail: View[] = $state([]);
	let current: View | null = $state(
		initial ? { id: recipeId, props: initial.recipe, canEdit: initial.canEdit } : null
	);
	let isLoading = $state(initial === undefined);
	let error: string | null = $state(null);
	let editorMode = $state(false);
	let confirmingDelete = $state(false);
	let busy = $state(false);
	let isLiked = $state(liked);
	let draft = $state(toDraft(initial?.recipe));

	const rootActive = $derived(trail.length === 0);

	async function open(id: number) {
		isLoading = true;
		error = null;
		const result = await getRecipe(id);
		if (result.ok && result.data.recipe?.Id !== undefined) {
			if (current) trail = [...trail, current];
			current = {
				id: result.data.recipe.Id,
				props: result.data.recipe.Props ?? {},
				canEdit: entryCanEdit(result.data)
			};
		} else {
			error = errorMessage(result.ok ? 404 : result.status);
		}
		isLoading = false;
	}

	onMount(() => {
		if (initial === undefined) open(recipeId);
	});

	function goBack() {
		current = trail[trail.length - 1] ?? null;
		trail = trail.slice(0, -1);
		editorMode = false;
		confirmingDelete = false;
		error = null;
	}

	function toggleEdit() {
		if (!editorMode && current) draft = toDraft(current.props);
		editorMode = !editorMode;
		error = null;
	}

	async function handleSave() {
		if (!current || busy || draft.name.trim() === '') return;
		busy = true;
		error = null;
		const updated = { ...current.props, ...cleanDraft(draft) };
		const result = await updateRecipe(current.id, updated);
		if (result.ok) {
			current = { ...current, props: updated };
			editorMode = false;
			onSaved?.();
		} else {
			error = errorMessage(result.status);
		}
		busy = false;
	}

	async function handleDelete() {
		if (!current || busy) return;
		busy = true;
		error = null;
		const result = await deleteRecipe(current.id);
		busy = false;
		if (!result.ok && result.status !== 404) {
			error = errorMessage(result.status);
			confirmingDelete = false;
			return;
		}
		confirmingDelete = false;
		onSaved?.();
		if (rootActive) {
			closeModal();
		} else {
			goBack();
		}
	}

	async function toggleLike() {
		if (busy) return;
		busy = true;
		error = null;
		const result = isLiked ? await unlikeRecipe(recipeId) : await likeRecipe(recipeId);
		if (result.ok) {
			isLiked = !isLiked;
			onLikeToggle?.(isLiked);
		} else {
			error = errorMessage(result.status);
		}
		busy = false;
	}
</script>

<div class="modal modal-open" transition:fade>
	<div class="modal-box max-h-screen w-full max-w-3xl overflow-y-auto">
		<div class="bg-base-100 sticky top-0 z-7">
			<div class="flex items-center justify-between p-2">
				<div class="flex items-center gap-2">
					<h3 class="text-lg font-bold">{current?.props.name ?? recipe()}</h3>
					{#if onLikeToggle && rootActive}
						<button
							class="btn btn-ghost btn-sm"
							class:text-error={isLiked}
							onclick={toggleLike}
							disabled={busy}
							title={favourite()}
						>
							{#if isLiked}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="currentColor"
									class="h-5 w-5"
								>
									<path
										d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3A5.5 5.5 0 0112 5.052 5.5 5.5 0 0116.313 3c2.973 0 5.437 2.322 5.437 5.25 0 3.925-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001a.752.752 0 01-.704 0l-.003-.001z"
									/>
								</svg>
							{:else}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="1.5"
									stroke="currentColor"
									class="h-5 w-5"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"
									/>
								</svg>
							{/if}
						</button>
					{/if}
				</div>
				<div class="space-x-2">
					{#if !rootActive && !editorMode}
						<button class="btn btn-ghost btn-sm" onclick={goBack}>{back()}</button>
					{/if}
					{#if current?.canEdit}
						<button class="btn btn-secondary btn-sm" onclick={toggleEdit}>
							{editorMode ? back() : edit()}
						</button>
						{#if editorMode}
							<button
								class="btn btn-accent btn-sm"
								onclick={handleSave}
								disabled={busy || draft.name.trim() === ''}
							>
								{save()}
							</button>
						{:else}
							<button
								class="btn btn-error btn-outline btn-sm"
								onclick={() => (confirmingDelete = true)}
								disabled={busy || confirmingDelete}
							>
								{delete_recipe()}
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

		{#if error}
			<div class="alert alert-error alert-soft mx-2 mb-2 py-2 text-sm" role="alert">{error}</div>
		{/if}

		{#if confirmingDelete}
			<div class="alert alert-warning mx-2 mb-2 flex items-center justify-between" role="alert">
				<span>{delete_recipe_confirm()}</span>
				<div class="flex gap-2">
					<button class="btn btn-ghost btn-sm" onclick={() => (confirmingDelete = false)}>
						{cancel()}
					</button>
					<button class="btn btn-error btn-sm" onclick={handleDelete} disabled={busy}>
						{delete_action()}
					</button>
				</div>
			</div>
		{/if}

		{#if isLoading}
			<div class="flex justify-center p-8">
				<span class="loading loading-spinner loading-lg" aria-label={loading()}></span>
			</div>
		{:else if current}
			<div class="flex flex-col gap-4 p-2">
				{#if editorMode}
					<RecipeFields bind:draft />
				{:else}
					<div>
						<span class="label font-semibold">{recipe()}</span>
						<p class="text-lg">{current.props.name ?? '-'}</p>
					</div>

					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<div>
							<span class="label font-semibold">{origin()}</span>
							<p>{current.props.origin ?? '-'}</p>
						</div>
						<div>
							<span class="label font-semibold">{category()}</span>
							<p>{current.props.category ?? '-'}</p>
						</div>
					</div>

					<div>
						<span class="label font-semibold">{description()}</span>
						<p>{current.props.description ?? '-'}</p>
					</div>

					<div>
						<span class="label font-semibold">{ingredients()}</span>
						{#if current.props.ingredients && current.props.ingredients.length > 0}
							<ul class="list-disc pl-5">
								{#each current.props.ingredients as item}
									<li>{item}</li>
								{/each}
							</ul>
						{:else}
							<p>-</p>
						{/if}
					</div>

					<div>
						<span class="label font-semibold">{instructions()}</span>
						{#if current.props.instructions && current.props.instructions.length > 0}
							<ol class="list-decimal pl-5">
								{#each current.props.instructions as step}
									<li class="mb-1">{step}</li>
								{/each}
							</ol>
						{:else}
							<p>-</p>
						{/if}
					</div>

					<div>
						<span class="label font-semibold">{notes()}</span>
						<p>{current.props.notes ?? '-'}</p>
					</div>

					{#key current.id}
						<div class="divider my-0"></div>
						<RecipeVariations recipeId={current.id} baseProps={current.props} onOpen={open} />
						<div class="divider my-0"></div>
						<RecipeComments recipeId={current.id} {currentUserId} />
					{/key}
				{/if}
			</div>
		{/if}
	</div>
</div>
