<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { setFamilyTreeNodes } from './setFamilyTreeNodes';
	import { user } from '../stores';
	import { onMount } from 'svelte';
	export let id = '';

	let auth_token = '';

	user.subscribe((value) => {
		if (value) {
			auth_token = value.access_token;
		}
	});

	let relationship = '';
	let dialog: HTMLDialogElement;

	let payload = {
		relationship: {
			first_person_id: id,
			relationship: relationship,
			direction: '->'
		},
		person: {
			first_name: '',
			middle_name: '',
			last_name: '',
			mothers_first_name: '',
			mothers_last_name: '',
			born: '',
			birthplace: '',
			titles: [],
			residence: '',
			death: '',
			deathplace: '',
			life_events: [],
			occupation: [],
			occupation_to_display: '',
			others_said: {},
			photos: {},
			profile_picture: 'https://cdn-icons-png.flaticon.com/512/3607/3607444.png'
		}
	};

	function handleSubmit(event: Event) {
		event.preventDefault();

		if (
			id &&
			relationship &&
			payload.person.first_name &&
			payload.person.last_name &&
			payload.person.born &&
			payload.person.mothers_first_name &&
			payload.person.mothers_last_name
		) {
			fetch(PUBLIC_API_URL + '/createRelationshipAndPerson', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: auth_token
				},
				body: JSON.stringify(payload)
			}).then((response) => {
				if (response.ok) {
					setFamilyTreeNodes();
					dialog.close();
				} else {
					alert('Failed to add family member! with status ' + response.status);
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
		<h1>Create a family members profile</h1>
		<p>
			You can add a family member to your family tree by filling out the form below. You can update
			the information later.
		</p>
		<form on:submit={handleSubmit} method="dialog">
			<label for="id">ID:</label>
			<input type="text" id="id" bind:value={id} required />

			<label for="relationship">Relationship type:</label>
			<select id="relationship" bind:value={relationship} required>
				<option value="" disabled selected>Choose one ...</option>
				<option value="Parent">Parent</option>
				<option value="Child">Child</option>
				<option value="Spouse">Spouse</option>
				<option value="Sibling">Sibling</option>
			</select>

			<label for="firstName">First Name:</label>
			<input type="text" id="firstName" bind:value={payload.person.first_name} required />

			<label for="middleName">Middle Name:</label>
			<input type="text" id="middleName" bind:value={payload.person.middle_name} />

			<label for="lastName">Last Name:</label>
			<input type="text" id="lastName" bind:value={payload.person.last_name} required />

			<label for="born">Born:</label>
			<input type="date" id="born" bind:value={payload.person.born} required />

			<label for="birthplace">Birthplace:</label>
			<input type="text" id="birthplace" bind:value={payload.person.birthplace} required />

			<label for="mothersFirstName">Mother's First Name:</label>
			<input
				type="text"
				id="mothersFirstName"
				bind:value={payload.person.mothers_first_name}
				required
			/>

			<label for="mothersLastName">Mother's Last Name:</label>
			<input
				type="text"
				id="mothersLastName"
				bind:value={payload.person.mothers_last_name}
				required
			/>

			<label for="death">Death:</label>
			<input type="date" id="death" bind:value={payload.person.death} />

			<label for="deathplace">Place of death:</label>
			<input type="text" id="deathplace" bind:value={payload.person.deathplace} />

			<label for="residence">Residence:</label>
			<input type="text" id="residence" bind:value={payload.person.residence} />

			<label for="titles">Titles:</label>
			<input type="text" id="titles" bind:value={payload.person.titles} />

			<button type="submit" class="btn btn-primary">Add</button>
		</form>
	</div>
</dialog>
