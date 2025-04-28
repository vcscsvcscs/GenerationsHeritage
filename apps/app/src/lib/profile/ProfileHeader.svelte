<script lang="ts">
	import { onMount } from 'svelte';
	import type { components } from '$lib/api/api.gen';
	import {
		male,
		female,
		intersex,
		other,
		change_profile_picture,
		biological_sex,
		born,
		died,
		email,
		first_name,
		id,
		last_name,
		middle_name,
		mothers_first_name,
		mothers_last_name,
		profile_picture
	} from '$lib/paraglide/messages';
	import { callMessageFunction } from '$lib/i18n';
	import type { MessageKeys } from '$lib/i18n';

	export let person: components['schemas']['PersonProperties'] & {
		id?: string;
	};
	export let editorMode = false;
	export let onChange: (field: keyof components['schemas']['PersonProperties'], value: any) => void;

	let birth_date: HTMLInputElement;
	let death_date: HTMLInputElement;
	onMount(() => {
		if (birth_date) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: birth_date,
					onSelect: function (date) {
						birth_date.value = date.toISOString().split('T')[0];
						onChange('born', date.toISOString().split('T')[0]);
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
		if (death_date) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: death_date,
					onSelect: function (date) {
						death_date.value = date.toISOString().split('T')[0];
						onChange('died', date.toISOString().split('T')[0]);
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
	});
</script>

<div class="flex flex-col gap-6 md:flex-row">
	<div class="flex flex-shrink-0 flex-col items-center gap-2">
		<img
			src={person.profile_picture || 'https://cdn-icons-png.flaticon.com/512/10628/10628885.png'}
			alt={profile_picture()}
			class="h-48 w-48 rounded-lg object-cover shadow-md"
		/>
		{#if editorMode}
			<button class="btn bg-neutral text-neutral-content btn-xs" on:click={() => {}}>
				{change_profile_picture()}
			</button>
		{/if}
	</div>
	<div class="grid flex-1 grid-cols-1 gap-4 md:grid-cols-2">
		<div>
			<p>
				<strong>{first_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.first_name}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.first_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{last_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.last_name}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.last_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{middle_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.middle_name}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.middle_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{born()}:</strong>
				{#if editorMode}<input
					type="text"
					class="pika-single w-full"
					id="birth_date"
					bind:this={birth_date}
					bind:value={person.born}
				/>
				{:else}{person.born ?? '-'}{/if}
			</p>
			<p>
				<strong>{died()}:</strong>
				{#if editorMode}<input
					type="text"
					class="pika-single w-full"
					id="death_date"
					placeholder={died()}
					bind:this={death_date}
					bind:value={person.died}
				/>{:else}{person.died ?? '-'}{/if}
			</p>
			<p>
				<strong>{biological_sex()}:</strong>
				{#if editorMode}
					<select
						name="biological_sex"
						class="select select-bordered select-sm w-full"
						id="biological_sex"
						bind:value={person.biological_sex}
						placeholder={biological_sex()}
					>
						<option value="male">{male()} </option>
						<option value="female">{female()} </option>
						<option value="intersex">{intersex()} </option>
						<option value="other">{other()} </option>
					</select>
				{:else}{callMessageFunction(person.biological_sex as MessageKeys) ?? '-'}{/if}
			</p>
		</div>
		<div>
			<p>
				<strong>{email()}:</strong>
				{#if editorMode}<input
						bind:value={person.email}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.email ?? '-'}{/if}
			</p>
			<p>
				<strong>{mothers_first_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.mothers_first_name}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.mothers_first_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{mothers_last_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.mothers_last_name}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.mothers_last_name ?? '-'}{/if}
			</p>
			<p><strong> {id()}:</strong>{person.id ?? '-'}</p>
			<p><strong> Limit:</strong>{person.limit ?? '-'}</p>
		</div>
	</div>
</div>
