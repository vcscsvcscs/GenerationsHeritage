<script lang="ts">
	import { callMessageFunction } from '$lib/i18n';
	import type { MessageKeys } from '$lib/i18n';
	import type { components } from '$lib/api/api.gen';

	export let person: components['schemas']['PersonProperties'];
	export let editorMode = false;
	export let onChange: (field: keyof components['schemas']['PersonProperties'], value: any) => void;
	const skipFields = [
		'id',
		'first_name',
		'last_name',
		'born',
		'died',
		'middle_name',
		'biological_sex',
		'email',
		'limit',
		'mothers_first_name',
		'mothers_last_name',
		'profile_picture',
		'photos',
		'videos',
		'life_events',
		'residence',
		'medications',
		'medical_conditions',
		'languages',
		'notes',
		'audios',
		'google_id'
	];
</script>

<div class="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2">
	{#each Object.entries(person) as [key, value]}
		{#if !skipFields.includes(key) && ((value !== undefined && value !== null) || editorMode)}
			<div>
				<label class="label font-semibold"
					>{callMessageFunction(key as MessageKeys) || key}:
					{#if editorMode}
						{#if typeof value === 'string'}
							{#if value.length > 100}
								<textarea
									bind:value={person[key as keyof components['schemas']['PersonProperties']]}
									class="textarea textarea-bordered textarea-sm w-full"
									oninput={(e) =>
										onChange(
											key as keyof components['schemas']['PersonProperties'],
											String(person[key as keyof components['schemas']['PersonProperties']])
										)}
								></textarea>
							{:else}
								<input
									type="text"
									class="input input-bordered input-sm w-full"
									bind:value={person[key as keyof components['schemas']['PersonProperties']]}
									oninput={() =>
										onChange(
											key as keyof components['schemas']['PersonProperties'],
											String(person[key as keyof components['schemas']['PersonProperties']])
										)}
								/>
							{/if}
						{:else if typeof value === 'boolean'}
							<input
								type="checkbox"
								class="checkbox checkbox-primary"
								bind:value={person[key as keyof components['schemas']['PersonProperties']]}
								onchange={(e) =>
									onChange(
										key as keyof components['schemas']['PersonProperties'],
										Boolean(person[key as keyof components['schemas']['PersonProperties']])
									)}
							/>
						{:else if typeof value === 'number'}
							<input
								type="number"
								class="input input-bordered input-sm w-full"
								bind:value={person[key as keyof components['schemas']['PersonProperties']]}
								oninput={(e) =>
									onChange(
										key as keyof components['schemas']['PersonProperties'],
										Number(person[key as keyof components['schemas']['PersonProperties']])
									)}
							/>
						{/if}
					{:else}
						<p>{value ?? '-'}</p>
					{/if}
				</label>
			</div>
		{/if}
	{/each}
</div>
