<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { close, cookbook, recipe, no_recipes } from '$lib/paraglide/messages';
	import RecipeModal from './RecipeModal.svelte';

	let {
		closeModal
	}: {
		closeModal: () => void;
	} = $props();

	let entries: Array<{
		recipe: { Id: number; Props: Record<string, unknown> };
		added_by: { id: number; first_name?: string; last_name?: string; profile_picture?: string };
		relationship: Record<string, unknown>;
	}> = $state([]);
	let likedRecipeIds: Set<number> = $state(new Set());
	let isLoading = $state(true);
	let selectedEntry: { id: number; props: Record<string, unknown> } | undefined = $state(undefined);

	onMount(() => {
		fetchCookbook();
		fetchMyLikedRecipes();
	});

	async function fetchCookbook() {
		isLoading = true;
		try {
			const response = await fetch('/api/cookbook?distance=5');
			if (response.ok) {
				const data = await response.json();
				entries = data?.entries ?? [];
			}
		} catch (e) {
			console.error('Error fetching cookbook:', e);
		} finally {
			isLoading = false;
		}
	}

	async function fetchMyLikedRecipes() {
		try {
			const response = await fetch('/api/recipe/person/me');
			if (response.ok) {
				const data = (await response.json()) as { recipes?: Array<{ Id: number }> };
				likedRecipeIds = new Set((data?.recipes ?? []).map((r) => r.Id));
			}
		} catch (e) {
			console.error('Error fetching liked recipes:', e);
		}
	}

	function getPersonName(person: { first_name?: string; last_name?: string }): string {
		return [person.first_name, person.last_name].filter(Boolean).join(' ') || '?';
	}

	function handleLikeToggle(recipeId: number, liked: boolean) {
		if (liked) {
			likedRecipeIds.add(recipeId);
		} else {
			likedRecipeIds.delete(recipeId);
		}
		likedRecipeIds = new Set(likedRecipeIds);
	}
</script>

{#if selectedEntry}
	<RecipeModal
		recipeData={selectedEntry.props}
		recipeId={selectedEntry.id}
		editable={false}
		liked={likedRecipeIds.has(selectedEntry.id)}
		closeModal={() => {
			selectedEntry = undefined;
		}}
		onLikeToggle={(liked) => handleLikeToggle(selectedEntry!.id, liked)}
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

			{#if isLoading}
				<div class="flex justify-center p-8">
					<span class="loading loading-spinner loading-lg"></span>
				</div>
			{:else if entries.length === 0}
				<p class="text-base-content/60 p-8 text-center">{no_recipes()}</p>
			{:else}
				<div class="flex flex-col gap-2 p-2">
					{#each entries as entry}
						<button
							class="card bg-base-200 hover:bg-base-300 w-full cursor-pointer text-left shadow-sm transition-colors"
							onclick={() => {
								selectedEntry = { id: entry.recipe.Id, props: entry.recipe.Props };
							}}
						>
							<div class="card-body p-4">
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2">
										<h4 class="card-title text-base">
											{entry.recipe.Props?.name ?? recipe()}
										</h4>
										{#if likedRecipeIds.has(entry.recipe.Id)}
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
										<span>{getPersonName(entry.added_by)}</span>
									</div>
								</div>
								<div class="text-base-content/60 flex gap-3 text-sm">
									{#if entry.recipe.Props?.category}
										<span class="badge badge-outline badge-sm">
											{entry.recipe.Props.category}
										</span>
									{/if}
									{#if entry.recipe.Props?.origin}
										<span>{entry.recipe.Props.origin}</span>
									{/if}
								</div>
								{#if entry.recipe.Props?.description}
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
