<script lang="ts">
	import type { PageData } from './$types';
	import {
		register,
		title,
		family_tree,
		welcome,
		site_intro,
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
		intersex,
		invite_code,
		have_invite_code
	} from '$lib/paraglide/messages';
	import { onMount } from 'svelte';
	import { enhance } from '$app/forms';
	import FamilyTree from '../../highresolution_icon_no_background_croped.png';
	let {
		data,
		form
	}: {
		data: PageData;
		form: {
			message: string;
		};
	} = $props();

	let birth_date: HTMLInputElement;
	let birth_date_value: HTMLInputElement;
	onMount(() => {
		if (birth_date) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: birth_date,
					onOpen: function () {
						birth_date_value.placeholder = '';
					},
					onSelect: function (date) {
						birth_date_value.value = date.toISOString().split('T')[0];
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
	});

	let showInviteInput = $state(false);
	function toggleInviteInput() {
		showInviteInput = !showInviteInput;
	}
</script>

<svelte:head>
	<title>{title({ page: register() })}</title>
</svelte:head>

<div class="hero bg-base-200 min-h-screen">
	<div class="hero-content flex-col lg:flex-row-reverse">
		<div class="max-w-xxl flex flex-col items-center justify-center text-center">
			<figure class="top-margin-10 max-w-sm px-10 pt-10">
				<img src={FamilyTree} alt={family_tree()} class="rounded-xl" />
			</figure>
			<h1 class="text-5xl font-bold">{welcome()}</h1>
			<p class="py-6">
				{site_intro()}
			</p>
		</div>
		<div class="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
			<div class="card-body">
				<form method="POST" action="?/register" use:enhance>
					<fieldset class="fieldset">
						{#if form?.message}
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
								<span>{form.message}</span>
							</div>
						{/if}
						<label class="fieldset-label" for="email">{email()}</label>
						<input
							type="email"
							name="email"
							class="input"
							placeholder={email()}
							value={data.props?.email}
						/>
						<input
							type="text"
							name="google_id"
							class="hidden"
							id="google_id"
							placeholder="Google ID"
							value={data.props?.google_id}
						/>
						<label class="fieldset-label" for="first_name">{first_name()}</label>
						<input
							type="text"
							class="input"
							name="first_name"
							id="first_name"
							placeholder={first_name()}
							value={data.props?.first_name}
						/>
						<label class="fieldset-label" for="last_name">{last_name()}</label>
						<input
							type="text"
							class="input"
							name="last_name"
							id="last_name"
							placeholder={last_name()}
							value={data.props?.last_name}
						/>
						<label class="fieldset-label" for="birth_date">{born()}</label>
						<input
							type="text"
							class="input pika-single"
							id="birth_date"
							placeholder={born()}
							bind:this={birth_date}
						/>
						<input type="text" class="hidden" name="birth_date" bind:this={birth_date_value} />
						<label class="fieldset-label" for="biological_sex">{biological_sex()}</label>
						<select
							name="biological_sex"
							class="select select-bordered w-full max-w-xs"
							id="biological_sex"
							placeholder={biological_sex()}
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
						/>
						<label class="fieldset-label" for="mothers_first_name">{mothers_first_name()}</label>
						<input
							type="text"
							class="input"
							name="mothers_first_name"
							id="mothers_first_name"
							placeholder={mothers_first_name()}
						/>
						<button type="button" class="btn btn-soft mt-4 max-w-xs" onclick={toggleInviteInput}>
							{have_invite_code()}
						</button>
						{#if showInviteInput}
							<div class="mt-4">
								<label class="fieldset-label" for="invite_code">Meghívókód</label>
								<input
									type="text"
									class="input input-bordered w-full"
									name="invite_code"
									id="invite_code"
									placeholder={invite_code()}
								/>
							</div>
						{/if}

						<button class="btn btn-neutral mt-4 max-w-xs">{register()}</button>
					</fieldset>
				</form>
			</div>
		</div>
	</div>
</div>
