<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { onMount } from 'svelte';
	import { user } from '../stores';

	let auth_token = '';

	user.subscribe((value) => {
		if (value) {
			auth_token = value.access_token;
		}
	});

	let dialog: HTMLDialogElement;

	let payload = {
		first_name: '',
		middle_name: '',
		last_name: '',
		mothers_first_name: '',
		mothers_last_name: '',
		born: '',
		birthplace: '',
		titles: [],
		residence: '',
		life_events: [],
		occupation: [],
		occupation_to_display: '',
		photos: {},
		profile_picture: 'https://cdn-icons-png.flaticon.com/512/3607/3607444.png'
	};

	function handleSubmit(event: Event) {
		event.preventDefault();

		if (
			payload.first_name &&
			payload.last_name &&
			payload.born &&
			payload.mothers_first_name &&
			payload.mothers_last_name
		) {
			payload.born = new Date(payload.born).toISOString();
			fetch(PUBLIC_API_URL + '/person', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer' + auth_token
				},
				body: JSON.stringify(payload)
			}).then((response) => {
				if (response.ok) {
					dialog.close();
				} else {
					alert('Failed to create profile! with status ' + response.status);
				}
			});
		} else {
			alert('You have to fill out all the required fields!');
		}
	}

	onMount(() => {
		if (dialog) dialog.showModal();
	});
</script>

<dialog bind:this={dialog} class="modal bg-base-300">
	<div class="modal-box w-11/12 max-w-5xl flex flex-col gap-4">
		<h1 class="font-bold text-lg">Create your profile</h1>
		<p>To start building your family tree fill this form. You can add more information later.</p>
		<form on:submit={handleSubmit} method="dialog">
			<div class="grid grid-cols-2 gap-4">
				<label class="input input-bordered flex items-center gap-2" for="firstName">
					First Name:
					<input
						type="text"
						id="firstName"
						class="grow input-md"
						bind:value={payload.first_name}
						required
					/>
				</label>

				<label class="input input-bordered flex items-center gap-2" for="middleName">
					Middle Name:
					<input
						type="text"
						id="middleName"
						class="grow input-md"
						bind:value={payload.middle_name}
					/>
				</label>

				<label class="input input-bordered flex items-center gap-2" for="lastName">
					Last Name:
					<input
						type="text"
						id="lastName"
						class="grow input-md"
						bind:value={payload.last_name}
						required
					/>
				</label>

				<label class="input input-bordered flex items-center gap-2" for="born">
					Born:
					<input type="date" id="born" class="grow input-md" bind:value={payload.born} required />
				</label>
				<label class="input input-bordered flex items-center gap-2" for="birthplace">
					Birthplace:
					<input
						type="text"
						id="birthplace"
						class="grow input-md"
						bind:value={payload.birthplace}
						required
					/>
				</label>
				<label class="input input-bordered flex items-center gap-2" for="mothersFirstName">
					Mother's First Name:
					<input
						type="text"
						id="mothersFirstName"
						class="grow input-md"
						bind:value={payload.mothers_first_name}
						required
					/>
				</label>
				<label class="input input-bordered flex items-center gap-2" for="mothersLastName">
					Mother's Last Name:
					<input
						type="text"
						id="mothersLastName"
						class="grow input-md"
						bind:value={payload.mothers_last_name}
						required
					/>
				</label>
				<label class="input input-bordered flex items-center gap-2" for="residence">
					Residence:
					<input type="text" id="residence" class="grow input-md" bind:value={payload.residence} />
				</label>
				<label class="input input-bordered flex items-center gap-2" for="titles">
					Titles:
					<input type="text" id="titles" class="grow input-md" bind:value={payload.titles} />
				</label>
				<button type="submit" class="btn btn-primary w-40 place-self-end">Add</button>
			</div>
		</form>
	</div>
</dialog>
