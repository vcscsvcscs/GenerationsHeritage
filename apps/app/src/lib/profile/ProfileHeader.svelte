<script lang="ts">
	import { onMount } from 'svelte';
	import type { components } from '$lib/api/api.gen';
	import { v4 as uuidv4 } from 'uuid';
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
		profile_picture,
		create_invite_code,

		invite_code

	} from '$lib/paraglide/messages';
	import { callMessageFunction } from '$lib/i18n';
	import type { MessageKeys } from '$lib/i18n';

	export let person: components['schemas']['PersonProperties'] & {
		id?: string;
	};
	export let editorMode = false;
	export let onChange: (field: keyof components['schemas']['PersonProperties'], value: any) => void;
	let new_invite_code: string | undefined;

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
			<button class="btn btn-neutral btn-soft btn-xs" onclick={() => {}}>
				{change_profile_picture()}
			</button>
		{/if}
	</div>
	<div class="grid flex-1 grid-cols-1 gap-4 md:grid-cols-2">
		<div class="flex flex-col gap-2">
			<p>
				<strong>{first_name()}: </strong>
				{#if editorMode}<input
						bind:value={person.first_name}
						onchange={() => onChange('first_name', person.first_name)}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.first_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{last_name()}: </strong>
				{#if editorMode}<input
						bind:value={person.last_name}
						onchange={() => onChange('last_name', person.last_name)}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.last_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{middle_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.middle_name}
						onchange={() => onChange('middle_name', person.middle_name)}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.middle_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{born()}: </strong>
				{#if editorMode}<input
						type="text"
						class="pika-single w-full"
						id="birth_date"
						bind:this={birth_date}
						placeholder={person.born}
						onchange={() => onChange('born', birth_date.value)}
				/>
				{:else}{person.born ?? '-'}{/if}
			</p>
			<p>
				<strong>{died()}: </strong>
				{#if editorMode}<input
						type="text"
						class="pika-single w-full"
						id="death_date"
						placeholder={person.died??died()}
						bind:this={death_date}
						onchange={() => onChange('died', death_date.value)}
					/>{:else}{person.died ?? '-'}{/if}
			</p>
			<p>
				<strong>{biological_sex()}: </strong>
				{#if editorMode}
					<select
						name="biological_sex"
						class="select select-bordered select-sm w-full"
						id="biological_sex"
						bind:value={person.biological_sex}
						onchange={() => onChange('biological_sex', person.biological_sex)}
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
		<div class="flex flex-col gap-2">
			<p>
				<strong>{email()}:</strong>
				{#if editorMode}<input
						bind:value={person.email}
						onchange={() => onChange('email', person.email)}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.email ?? '-'}{/if}
			</p>
			<p>
				<strong>{mothers_first_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.mothers_first_name}
						onchange={() => onChange('mothers_first_name', person.mothers_first_name)}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.mothers_first_name ?? '-'}{/if}
			</p>
			<p>
				<strong>{mothers_last_name()}:</strong>
				{#if editorMode}<input
						bind:value={person.mothers_last_name}
						onchange={() => onChange('mothers_last_name', person.mothers_last_name)}
						class="input input-sm input-bordered w-full"
					/>{:else}{person.mothers_last_name ?? '-'}{/if}
			</p>
			<p><strong>{id()}: </strong>{' ' + (person.id ?? '-')}</p>
			<p><strong>Limit: </strong>{' ' + (person.limit ?? '-')}</p>
			{#if editorMode && (person.google_id === undefined || person.google_id === null || person.google_id === '')}
				{#if new_invite_code===undefined}
				<button class="btn btn-soft btn-accent btn-m" onclick={()=>{
					new_invite_code = uuidv4();
					person.invite_code = new_invite_code;
					onChange('invite_code',new_invite_code);
					}}>{create_invite_code()}</button>
				{:else}
				<p>
					<strong>{invite_code()}:</strong>{person.invite_code}
				</p>
				{/if}
			{/if}
		</div>
	</div>
</div>
