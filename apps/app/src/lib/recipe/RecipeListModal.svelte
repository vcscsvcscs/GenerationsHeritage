<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import {
		close,
		recipes,
		liked_recipes,
		new_recipe,
		no_recipes,
		save,
		cancel,
		recipe,
		description,
		origin,
		category,
		ingredients,
		instructions,
		notes,
		add,
		loading
	} from '$lib/paraglide/messages';
	import RecipeModal from './RecipeModal.svelte';

	let {
		personId,
		personName = '',
		useMyRecipes = false,
		closeModal
	}: {
		personId: number;
		personName?: string;
		useMyRecipes?: boolean;
		closeModal: () => void;
	} = $props();

	let recipeList: Array<{ Id: number; Props: any }> = $state([]);
	let isLoading = $state(true);
	let showCreateForm = $state(false);
	let selectedRecipe: { id: number; props: any } | undefined = $state(undefined);

	let newRecipe = $state({
		name: '',
		origin: '',
		category: '',
		description: '',
		ingredients: [''],
		instructions: [''],
		notes: ''
	});

	onMount(() => {
		fetchRecipes();
	});

	async function fetchRecipes() {
		isLoading = true;
		try {
			const endpoint = useMyRecipes ? '/api/recipe/person/me' : `/api/recipe/person/${personId}`;
			const response = await fetch(endpoint);
			if (response.ok) {
				const data = (await response.json()) as { recipes?: Array<{ Id: number; Props: any }> };
				recipeList = data?.recipes ?? [];
			}
		} catch (e) {
			console.error('Error fetching recipes:', e);
		} finally {
			isLoading = false;
		}
	}

	async function createRecipe() {
		const cleanRecipe = {
			...newRecipe,
			ingredients: newRecipe.ingredients.filter((i) => i.trim() !== ''),
			instructions: newRecipe.instructions.filter((i) => i.trim() !== '')
		};

		try {
			const response = await fetch(`/api/recipe/person/${personId}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					recipe: cleanRecipe,
					relationship: { favourite: true, like_it: true }
				})
			});
			if (response.ok) {
				showCreateForm = false;
				newRecipe = {
					name: '',
					origin: '',
					category: '',
					description: '',
					ingredients: [''],
					instructions: [''],
					notes: ''
				};
				await fetchRecipes();
			} else {
				alert('Error creating recipe');
			}
		} catch (e) {
			alert('Error creating recipe: ' + e);
		}
	}

	function addIngredient() {
		newRecipe.ingredients = [...newRecipe.ingredients, ''];
	}

	function addInstruction() {
		newRecipe.instructions = [...newRecipe.instructions, ''];
	}
</script>

{#if selectedRecipe}
	<RecipeModal
		recipeData={selectedRecipe.props}
		recipeId={selectedRecipe.id}
		editable={true}
		closeModal={() => {
			selectedRecipe = undefined;
		}}
		onSaved={fetchRecipes}
	/>
{:else}
	<div class="modal modal-open" transition:fade>
		<div class="modal-box max-h-screen w-full max-w-3xl overflow-y-auto">
			<div class="bg-base-100 z-7 sticky top-0">
				<div class="flex items-center justify-between p-2">
					<h3 class="text-lg font-bold">
						{useMyRecipes ? liked_recipes() : personName ? `${personName} - ${recipes()}` : recipes()}
					</h3>
					<div class="space-x-2">
						{#if !showCreateForm && !useMyRecipes}
							<button class="btn btn-accent btn-sm" onclick={() => (showCreateForm = true)}>
								{new_recipe()}
							</button>
						{/if}
						<button class="btn btn-error btn-sm" onclick={closeModal}>
							{close()}
						</button>
					</div>
				</div>
				<div class="divider"></div>
			</div>

			{#if showCreateForm}
				<!-- Create Recipe Form -->
				<div class="flex flex-col gap-3 p-2">
					<div>
						<label class="label font-semibold">{recipe()}</label>
						<input type="text" class="input input-bordered w-full" bind:value={newRecipe.name} />
					</div>
					<div class="grid grid-cols-1 gap-3 md:grid-cols-2">
						<div>
							<label class="label font-semibold">{origin()}</label>
							<input type="text" class="input input-bordered w-full" bind:value={newRecipe.origin} />
						</div>
						<div>
							<label class="label font-semibold">{category()}</label>
							<input
								type="text"
								class="input input-bordered w-full"
								bind:value={newRecipe.category}
							/>
						</div>
					</div>
					<div>
						<label class="label font-semibold">{description()}</label>
						<textarea
							class="textarea textarea-bordered w-full"
							rows="2"
							bind:value={newRecipe.description}
						></textarea>
					</div>
					<div>
						<label class="label font-semibold">{ingredients()}</label>
						{#each newRecipe.ingredients as item, i}
							<div class="flex items-center gap-2 mb-1">
								<input
									type="text"
									class="input input-bordered input-sm flex-1"
									bind:value={newRecipe.ingredients[i]}
								/>
								<button
									class="btn btn-xs btn-ghost text-error"
									onclick={() => {
										newRecipe.ingredients = newRecipe.ingredients.filter((_, idx) => idx !== i);
									}}
								>
									&#10005;
								</button>
							</div>
						{/each}
						<button class="btn btn-accent btn-xs mt-1" onclick={addIngredient}>
							{add()}
						</button>
					</div>
					<div>
						<label class="label font-semibold">{instructions()}</label>
						{#each newRecipe.instructions as step, i}
							<div class="flex items-center gap-2 mb-1">
								<span class="text-sm font-mono w-6">{i + 1}.</span>
								<textarea
									class="textarea textarea-bordered textarea-sm flex-1"
									bind:value={newRecipe.instructions[i]}
								></textarea>
								<button
									class="btn btn-xs btn-ghost text-error"
									onclick={() => {
										newRecipe.instructions = newRecipe.instructions.filter(
											(_, idx) => idx !== i
										);
									}}
								>
									&#10005;
								</button>
							</div>
						{/each}
						<button class="btn btn-accent btn-xs mt-1" onclick={addInstruction}>
							{add()}
						</button>
					</div>
					<div>
						<label class="label font-semibold">{notes()}</label>
						<textarea
							class="textarea textarea-bordered w-full"
							rows="2"
							bind:value={newRecipe.notes}
						></textarea>
					</div>
					<div class="flex gap-2 justify-end mt-2">
						<button class="btn btn-ghost btn-sm" onclick={() => (showCreateForm = false)}>
							{cancel()}
						</button>
						<button class="btn btn-primary btn-sm" onclick={createRecipe} disabled={!newRecipe.name}>
							{save()}
						</button>
					</div>
				</div>
			{:else if isLoading}
				<div class="flex justify-center p-8">
					<span class="loading loading-spinner loading-lg"></span>
				</div>
			{:else if recipeList.length === 0}
				<p class="text-center p-8 text-base-content/60">{no_recipes()}</p>
			{:else}
				<!-- Recipe List -->
				<div class="flex flex-col gap-2 p-2">
					{#each recipeList as r}
						<button
							class="card bg-base-200 shadow-sm hover:bg-base-300 transition-colors cursor-pointer w-full text-left"
							onclick={() => {
								selectedRecipe = { id: r.Id, props: r.Props };
							}}
						>
							<div class="card-body p-4">
								<h4 class="card-title text-base">{r.Props?.name ?? recipe()}</h4>
								<div class="flex gap-3 text-sm text-base-content/60">
									{#if r.Props?.category}
										<span class="badge badge-outline badge-sm">{r.Props.category}</span>
									{/if}
									{#if r.Props?.origin}
										<span>{r.Props.origin}</span>
									{/if}
								</div>
								{#if r.Props?.description}
									<p class="text-sm mt-1 line-clamp-2">{r.Props.description}</p>
								{/if}
							</div>
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}
