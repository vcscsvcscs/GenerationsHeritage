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
			fetch(PUBLIC_API_URL + '/person', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: auth_token
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
	<div class="modal-box w-11/12 max-w-5xl">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
		</form>
		<h1>Create your profile</h1>
		<p>
			Create your profile to start building your family tree. You can add more information later.
		</p>
		<form on:submit={handleSubmit} method="dialog">
			<label for="firstName">First Name:</label>
			<input type="text" id="firstName" bind:value={payload.first_name} required />

			<label for="middleName">Middle Name:</label>
			<input type="text" id="middleName" bind:value={payload.middle_name} />

			<label for="lastName">Last Name:</label>
			<input type="text" id="lastName" bind:value={payload.last_name} required />

			<label for="born">Born:</label>
			<input type="date" id="born" bind:value={payload.born} required />

			<label for="birthplace">Birthplace:</label>
			<input type="text" id="birthplace" bind:value={payload.birthplace} required />

			<label for="mothersFirstName">Mother's First Name:</label>
			<input type="text" id="mothersFirstName" bind:value={payload.mothers_first_name} required />

			<label for="mothersLastName">Mother's Last Name:</label>
			<input type="text" id="mothersLastName" bind:value={payload.mothers_last_name} required />

			<label for="residence">Residence:</label>
			<input type="text" id="residence" bind:value={payload.residence} />

			<label for="titles">Titles:</label>
			<input type="text" id="titles" bind:value={payload.titles} />

			<button type="submit" class="btn btn-primary">Add</button>
		</form>
	</div>
</dialog>
