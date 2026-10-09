<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { close, cookbook, recipe, no_recipes } from '$lib/paraglide/messages';
	import { getCookbook, getPersonRecipes } from './api';
	import { errorMessage } from './errors';
	import { canEdit, likedRecipeIds, personName } from './model';
	import type { CookbookEntry } from './model';
	import RecipeModal from './RecipeModal.svelte';

	let {
		currentUserId = null,
		closeModal
	}: {
		currentUserId?: number | null;
		closeModal: () => void;
	} = $props();

	let entries: CookbookEntry[] = $state([]);
	let likedIds: Set<number> = $state(new Set());
	let isLoading = $state(true);
	let error: string | null = $state(null);
	let selectedEntry: CookbookEntry | undefined = $state(undefined);

	onMount(() => {
		fetchCookbook();
		fetchMyLikedRecipes();
	});

	async function fetchCookbook() {
		isLoading = true;
		error = null;
		const result = await getCookbook(5);
		if (result.ok) {
			entries = (result.data.entries ?? []).filter((e) => e.recipe?.Id !== undefined);
		} else {
			error = errorMessage(result.status);
		}
		isLoading = false;
	}

	async function fetchMyLikedRecipes() {
		const result = await getPersonRecipes('me');
		if (result.ok) likedIds = likedRecipeIds(result.data.entries);
	}

	function handleLikeToggle(recipeId: number, liked: boolean) {
		const next = new Set(likedIds);
		if (liked) next.add(recipeId);
		else next.delete(recipeId);
		likedIds = next;
	}
</script>

{#if selectedEntry?.recipe?.Id !== undefined}
	{@const selectedId = selectedEntry.recipe.Id}
	<RecipeModal
		recipeId={selectedId}
		initial={{ recipe: selectedEntry.recipe.Props ?? {}, canEdit: canEdit(selectedEntry) }}
		{currentUserId}
		liked={likedIds.has(selectedId)}
		closeModal={() => {
			selectedEntry = undefined;
		}}
		onSaved={fetchCookbook}
		onLikeToggle={(liked) => handleLikeToggle(selectedId, liked)}
	/>
{:else}
	<div class="modal modal-open" transition:fade>
		<div class="modal-box max-h-screen w-full max-w-3xl overflow-y-auto">
			<div class="bg-base-100 sticky top-0 z-7">
				<div class="flex items-center justify-between p-2">
					<h3 class="text-lg font-bold">{cookbook()}</h3>
					<button class="btn btn-error btn-sm" onclick={closeModal}>
						{close()}
					</button>
				</div>
				<div class="divider"></div>
			</div>

			{#if error}
				<div class="alert alert-error alert-soft mx-2 mb-2 py-2 text-sm" role="alert">{error}</div>
			{/if}

			{#if isLoading}
				<div class="flex justify-center p-8">
					<span class="loading loading-spinner loading-lg"></span>
				</div>
			{:else if entries.length === 0}
				<p class="text-base-content/60 p-8 text-center">{no_recipes()}</p>
			{:else}
				<div class="flex flex-col gap-2 p-2">
					{#each entries as entry (entry.recipe?.Id)}
						<button
							class="card bg-base-200 hover:bg-base-300 w-full cursor-pointer text-left shadow-sm transition-colors"
							onclick={() => {
								selectedEntry = entry;
							}}
						>
							<div class="card-body p-4">
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2">
										<h4 class="card-title text-base">
											{entry.recipe?.Props?.name ?? recipe()}
										</h4>
										{#if entry.recipe?.Id !== undefined && likedIds.has(entry.recipe.Id)}
											<svg
												xmlns="http://www.w3.org/2000/svg"
												viewBox="0 0 24 24"
												fill="currentColor"
												class="text-error h-4 w-4"
											>
												<path
													d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3A5.5 5.5 0 0112 5.052 5.5 5.5 0 0116.313 3c2.973 0 5.437 2.322 5.437 5.25 0 3.925-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001a.752.752 0 01-.704 0l-.003-.001z"
												/>
											</svg>
										{/if}
									</div>
									<div class="text-base-content/60 flex items-center gap-2 text-sm">
										{#if entry.added_by?.profile_picture}
											<img
												src={entry.added_by.profile_picture}
												alt=""
												class="h-6 w-6 rounded-full"
											/>
										{/if}
										<span>{personName(entry.added_by)}</span>
									</div>
								</div>
								<div class="text-base-content/60 flex gap-3 text-sm">
									{#if entry.recipe?.Props?.category}
										<span class="badge badge-outline badge-sm">
											{entry.recipe.Props.category}
										</span>
									{/if}
									{#if entry.recipe?.Props?.origin}
										<span>{entry.recipe.Props.origin}</span>
									{/if}
								</div>
								{#if entry.recipe?.Props?.description}
									<p class="mt-1 line-clamp-2 text-sm">{entry.recipe.Props.description}</p>
								{/if}
							</div>
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}
