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
		add
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

	let recipeList: Array<{ Id: number; Props: Record<string, unknown> }> = $state([]);
	let isLoading = $state(true);
	let showCreateForm = $state(false);
	let selectedRecipe: { id: number; props: Record<string, unknown> } | undefined =
		$state(undefined);

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
				const data = (await response.json()) as {
					recipes?: Array<{ Id: number; Props: Record<string, unknown> }>;
				};
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
			<div class="bg-base-100 sticky top-0 z-7">
				<div class="flex items-center justify-between p-2">
					<h3 class="text-lg font-bold">
						{useMyRecipes
							? liked_recipes()
							: personName
								? `${personName} - ${recipes()}`
								: recipes()}
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
						<span class="label font-semibold">{recipe()}</span>
						<input type="text" class="input input-bordered w-full" bind:value={newRecipe.name} />
					</div>
					<div class="grid grid-cols-1 gap-3 md:grid-cols-2">
						<div>
							<span class="label font-semibold">{origin()}</span>
							<input
								type="text"
								class="input input-bordered w-full"
								bind:value={newRecipe.origin}
							/>
						</div>
						<div>
							<span class="label font-semibold">{category()}</span>
							<input
								type="text"
								class="input input-bordered w-full"
								bind:value={newRecipe.category}
							/>
						</div>
					</div>
					<div>
						<span class="label font-semibold">{description()}</span>
						<textarea
							class="textarea textarea-bordered w-full"
							rows="2"
							bind:value={newRecipe.description}
						></textarea>
					</div>
					<div>
						<span class="label font-semibold">{ingredients()}</span>
						{#each newRecipe.ingredients as _, i}
							<div class="mb-1 flex items-center gap-2">
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
						<span class="label font-semibold">{instructions()}</span>
						{#each newRecipe.instructions as __, i}
							<div class="mb-1 flex items-center gap-2">
								<span class="w-6 font-mono text-sm">{i + 1}.</span>
								<textarea
									class="textarea textarea-bordered textarea-sm flex-1"
									bind:value={newRecipe.instructions[i]}
								></textarea>
								<button
									class="btn btn-xs btn-ghost text-error"
									onclick={() => {
										newRecipe.instructions = newRecipe.instructions.filter((_, idx) => idx !== i);
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
						<span class="label font-semibold">{notes()}</span>
						<textarea
							class="textarea textarea-bordered w-full"
							rows="2"
							bind:value={newRecipe.notes}
						></textarea>
					</div>
					<div class="mt-2 flex justify-end gap-2">
						<button class="btn btn-ghost btn-sm" onclick={() => (showCreateForm = false)}>
							{cancel()}
						</button>
						<button
							class="btn btn-primary btn-sm"
							onclick={createRecipe}
							disabled={!newRecipe.name}
						>
							{save()}
						</button>
					</div>
				</div>
			{:else if isLoading}
				<div class="flex justify-center p-8">
					<span class="loading loading-spinner loading-lg"></span>
				</div>
			{:else if recipeList.length === 0}
				<p class="text-base-content/60 p-8 text-center">{no_recipes()}</p>
			{:else}
				<!-- Recipe List -->
				<div class="flex flex-col gap-2 p-2">
					{#each recipeList as r}
						<button
							class="card bg-base-200 hover:bg-base-300 w-full cursor-pointer text-left shadow-sm transition-colors"
							onclick={() => {
								selectedRecipe = { id: r.Id, props: r.Props };
							}}
						>
							<div class="card-body p-4">
								<h4 class="card-title text-base">{r.Props?.name ?? recipe()}</h4>
								<div class="text-base-content/60 flex gap-3 text-sm">
									{#if r.Props?.category}
										<span class="badge badge-outline badge-sm">{r.Props.category}</span>
									{/if}
									{#if r.Props?.origin}
										<span>{r.Props.origin}</span>
									{/if}
								</div>
								{#if r.Props?.description}
									<p class="mt-1 line-clamp-2 text-sm">{r.Props.description}</p>
								{/if}
							</div>
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}
