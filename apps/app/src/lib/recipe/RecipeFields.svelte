<script lang="ts">
	import {
		recipe,
		description,
		origin,
		category,
		ingredients,
		instructions,
		notes,
		add
	} from '$lib/paraglide/messages';
	import type { RecipeDraft } from './model';

	let { draft = $bindable() }: { draft: RecipeDraft } = $props();
</script>

<div class="flex flex-col gap-3">
	<div>
		<span class="label font-semibold">{recipe()}</span>
		<input type="text" class="input input-bordered w-full" bind:value={draft.name} />
	</div>
	<div class="grid grid-cols-1 gap-3 md:grid-cols-2">
		<div>
			<span class="label font-semibold">{origin()}</span>
			<input type="text" class="input input-bordered w-full" bind:value={draft.origin} />
		</div>
		<div>
			<span class="label font-semibold">{category()}</span>
			<input type="text" class="input input-bordered w-full" bind:value={draft.category} />
		</div>
	</div>
	<div>
		<span class="label font-semibold">{description()}</span>
		<textarea class="textarea textarea-bordered w-full" rows="2" bind:value={draft.description}
		></textarea>
	</div>
	<div>
		<span class="label font-semibold">{ingredients()}</span>
		{#each draft.ingredients as _, i}
			<div class="mb-1 flex items-center gap-2">
				<input
					type="text"
					class="input input-bordered input-sm flex-1"
					bind:value={draft.ingredients[i]}
				/>
				<button
					class="btn btn-xs btn-ghost text-error"
					onclick={() => (draft.ingredients = draft.ingredients.filter((_, idx) => idx !== i))}
				>
					&#10005;
				</button>
			</div>
		{/each}
		<button
			class="btn btn-accent btn-xs mt-1"
			onclick={() => (draft.ingredients = [...draft.ingredients, ''])}
		>
			{add()}
		</button>
	</div>
	<div>
		<span class="label font-semibold">{instructions()}</span>
		{#each draft.instructions as __, i}
			<div class="mb-1 flex items-center gap-2">
				<span class="w-6 font-mono text-sm">{i + 1}.</span>
				<textarea
					class="textarea textarea-bordered textarea-sm flex-1"
					bind:value={draft.instructions[i]}
				></textarea>
				<button
					class="btn btn-xs btn-ghost text-error"
					onclick={() => (draft.instructions = draft.instructions.filter((_, idx) => idx !== i))}
				>
					&#10005;
				</button>
			</div>
		{/each}
		<button
			class="btn btn-accent btn-xs mt-1"
			onclick={() => (draft.instructions = [...draft.instructions, ''])}
		>
			{add()}
		</button>
	</div>
	<div>
		<span class="label font-semibold">{notes()}</span>
		<textarea class="textarea textarea-bordered w-full" rows="2" bind:value={draft.notes}
		></textarea>
	</div>
</div>
