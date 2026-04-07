<script lang="ts">
	import { callMessageFunction } from '$lib/i18n';
	import type { MessageKeys } from '$lib/i18n';
	import { add_note, notes, theme } from '$lib/paraglide/messages';
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
		'phone',
		'audios',
		'google_id',
		'invite_code'
	];
	let newNote = {
		title: ' ',
		note: ''
	};
</script>

<div class="mt-5 flex flex-col items-center justify-center gap-2">
	{#each person.notes ?? [] as note, i}
		<div class="card bg-base-100 relative w-full max-w-xl shadow-sm">
			<div class="card-body w-full p-4">
				{#if editorMode}
					<input
						type="text"
						class="input input-bordered input-sm mb-2 w-full"
						placeholder={theme()}
						bind:value={note.title}
						oninput={() => onChange('notes', person.notes)}
					/>
					<textarea
						class="textarea textarea-bordered textarea-sm w-full"
						placeholder={notes()}
						bind:value={note.note}
						oninput={() => onChange('notes', person.notes)}
					></textarea>
					<button
						type="button"
						class="btn btn-xs btn-ghost text-error absolute top-2 right-2 ml-2"
						aria-label="Remove note"
						onclick={() => {
							person.notes = (person.notes ?? []).filter((_, idx) => idx !== i);
							onChange('notes', person.notes);
						}}
					>
						&#10005;
					</button>
				{:else}
					<h2 class="card-title">{note.title}</h2>
					<p class="text-sm text-gray-500">{note.date}</p>
					<p>{note.note}</p>
				{/if}
			</div>
		</div>
	{/each}
	{#if editorMode}
		<button
			class="btn btn-accent btn-sm w-auto self-start"
			onclick={() => {
				const now = new Date();
				const formattedDate =
					now.getFullYear() +
					'-' +
					String(now.getMonth() + 1).padStart(2, '0') +
					'-' +
					String(now.getDate()).padStart(2, '0');
				person.notes = [...(person.notes ?? []), { title: '', note: '', date: formattedDate }];
				onChange('notes', person.notes);
			}}
		>
			{add_note()}
		</button>
	{/if}
</div>
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
