<script lang="ts">
	import { fade } from 'svelte/transition';
	import {
		register,
		close,
		born,
		mothers_first_name,
		mothers_last_name,
		last_name,
		first_name,
		email,
		biological_sex,
		male,
		female,
		other,
		intersex
	} from '$lib/paraglide/messages';
	import { onMount } from 'svelte';
	import type { components } from '$lib/api/api.gen.js';
	import { validatePersonRegistration, validateFamilyRelationship } from './validate_fields';

	let {
		closeModal = () => {},
		relationship = null,
	}: { closeModal : ()=>void,relationship: number | null } = $props();
	let birth_date: HTMLInputElement;
	let draftRelationship: components['schemas']['FamilyRelationship'] & {type: string} | null = {} as components['schemas']['FamilyRelationship'] & {type: string} | null;
	let draftPerson :components['schemas']['PersonRegistration'] = {} as components['schemas']['PersonRegistration'];
	let error: string | undefined | null = $state();

	function onClose() {
		closeModal();
	}

	async function create(event: SubmitEvent) {
		event.preventDefault();
		error = validatePersonRegistration(draftPerson);
		if (error) {
			return;
		}

		if (relationship !== null && draftRelationship !== null) {
			error = validateFamilyRelationship(draftRelationship);
			if (error) {
				return;
			}
		}

		const url = `/api/person`;
		const response = await fetch(url);
		result = await response.text();
		closeModal();
	}

	onMount(() => {
		if (birth_date) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: birth_date,
					onOpen: function () {
						birth_date.placeholder = '';
					},
					onSelect: function (date) {
						birth_date.value = date.toISOString();
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
	});
</script>

<div class="modal modal-open" transition:fade>
	<div class="modal-box max-h-screen w-full max-w-5xl overflow-y-auto">
		<div class="bg-base-100 sticky top-0 z-10">
			<button class="btn btn-error btn-sm" onclick={onClose}>
				{close()}
			</button>
			<div class="divider"></div>
		</div>
		<form onsubmit={create}>
			<fieldset class="fieldset">
				{#if error}
					<div role="alert" class="alert alert-error">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 shrink-0 stroke-current"
							fill="none"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
						<span>{error}</span>
					</div>
				{/if}
				{#if relationship !== undefined}
					<input type="hidden" name="relationship" value={relationship} />
				{/if}
				<label class="fieldset-label" for="email">{email()}</label>
				<input
					type="email"
					name="email"
					class="input"
					placeholder={email()}
					bind:value={draftPerson.email}
				/>
				<label class="fieldset-label" for="first_name">{first_name()}</label>
				<input
					type="text"
					class="input"
					name="first_name"
					id="first_name"
					placeholder={first_name()}
				/>
				<label class="fieldset-label" for="last_name">{last_name()}</label>
				<input
					type="text"
					class="input"
					name="last_name"
					id="last_name"
					placeholder={last_name()}
				/>
				<label class="fieldset-label" for="birth_date">{born()}</label>
				<input
					type="text"
					class="input pika-single"
					id="birth_date"
					placeholder={born()}
					bind:value={draftPerson.born}
					bind:this={birth_date}
				/>
				<label class="fieldset-label" for="biological_sex">{biological_sex()}</label>
				<select
					name="biological_sex"
					class="select select-bordered w-full max-w-xs"
					id="biological_sex"
					placeholder={biological_sex()}
					bind:value={draftPerson.biological_sex}
				>
					<option value="male">{male()} </option>
					<option value="female">{female()} </option>
					<option value="intersex">{intersex()} </option>
					<option value="other">{other()} </option>
				</select>
				<label class="fieldset-label" for="mothers_last_name">{mothers_last_name()}</label>
				<input
					type="text"
					class="input"
					name="mothers_last_name"
					id="mothers_last_name"
					placeholder={mothers_last_name()}
					bind:value={draftPerson.mothers_last_name}
				/>
				<label class="fieldset-label" for="mothers_first_name">{mothers_first_name()}</label>
				<input
					type="text"
					class="input"
					name="mothers_first_name"
					id="mothers_first_name"
					placeholder={mothers_first_name()}
					bind:value={draftPerson.mothers_first_name}
				/>
				<button class="btn btn-neutral mt-4">{register()}</button>
			</fieldset>
		</form>
	</div>
</div>
