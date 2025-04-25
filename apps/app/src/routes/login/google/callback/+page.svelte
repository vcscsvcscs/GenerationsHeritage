<script lang="ts">
	import type { PageProps } from './$types';
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
		allow_family_tree_admin_access
	} from '$lib/paraglide/messages';
	import FamilyTree from '../../highresolution_icon_no_background_croped.png';
	let { data, form }: PageProps = $props();
	import Pikaday from 'pikaday';

	let birth_date: HTMLInputElement;
	$effect(() => {
		if (birth_date) {
			const picker = new Pikaday({
				field: birth_date
			});
			return () => picker.destroy();
		}
	});
</script>

<svelte:head>
	<title>{title({ page: register() })}</title>
</svelte:head>

<div class="hero bg-base-200 min-h-screen">
	<div class="hero-content flex-col lg:flex-row-reverse">
		<div class="text-center lg:text-left">
			<figure class="top-margin-10 px-10 pt-10">
				<img src={FamilyTree} alt={family_tree()} class="rounded-xl" />
			</figure>
			<h1 class="text-5xl font-bold">{welcome()}</h1>
			<p class="py-6">
				{site_intro()}
			</p>
		</div>
		<div class="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
			<div class="card-body">
				<form method="POST" action="/register">
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
						<label class="fieldset-label" for="email">Email</label>
						<input type="email" class="input" placeholder="Email" value={data.props?.email} />
						<input
							type="text"
							class="hidden"
							id="google_id"
							placeholder="Google ID"
							value={data.props?.google_id}
						/>
						<label class="fieldset-label" for="first_name">{first_name()}</label>
						<input
							type="text"
							class="input"
							id="first_name"
							placeholder={first_name()}
							value={data.props?.first_name}
						/>
						<label class="fieldset-label" for="last_name">{last_name()}</label>
						<input
							type="text"
							class="input"
							id="last_name"
							placeholder={last_name()}
							value={data.props?.last_name}
						/>
						<label class="fieldset-label" for="birth_date">{born()}</label>
						<input
							type="text"
							class="input pika-single"
							id="birth_date"
							bind:this={birth_date}
							value={born()}
						/>
						<label class="fieldset-label" for="mothers_last_name">{mothers_last_name()}</label>
						<input
							type="text"
							class="input"
							id="mothers_last_name"
							placeholder={mothers_last_name()}
						/>
						<label class="fieldset-label" for="mothers_first_name">{mothers_first_name()}</label>
						<input
							type="text"
							class="input"
							id="mothers_first_name"
							placeholder={mothers_first_name()}
						/>
						<button class="btn btn-neutral mt-4">{register()}</button>
					</fieldset>
				</form>
			</div>
		</div>
	</div>
</div>
