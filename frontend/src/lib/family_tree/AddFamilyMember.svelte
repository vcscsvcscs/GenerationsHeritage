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
		<h1 class="font-bold text-lg">Create a family members profile</h1>
		<p>
			You can add a family member to your family tree by filling out the form below. You can update
			the information later.
		</p>
		<form on:submit={handleSubmit} method="dialog" class="grid grid-cols-2 gap-4">
			<label for="id" class="input input-bordered flex items-center gap-2">
				ID:
				<input type="text" id="id" class="grow input-md" bind:value={id} required />
			</label>
			<label for="relationship" class="input input-bordered flex items-center gap-2">
				Relationship type:
				<select id="relationship" bind:value={relationship} required>
					<option value="" disabled selected>Choose one ...</option>
					<option value="Parent">Parent</option>
					<option value="Child">Child</option>
					<option value="Spouse">Spouse</option>
					<option value="Sibling">Sibling</option>
				</select>
			</label>

			<label for="firstName" class="input input-bordered flex items-center gap-2">
				First Name:
				<input
					type="text"
					id="firstName"
					class="grow input-md"
					bind:value={payload.person.first_name}
					required
				/>
			</label>

			<label for="middleName" class="input input-bordered flex items-center gap-2">
				Middle Name:
				<input
					type="text"
					id="middleName"
					class="grow input-md"
					bind:value={payload.person.middle_name}
				/>
			</label>

			<label for="lastName" class="input input-bordered flex items-center gap-2">
				Last Name:
				<input
					type="text"
					id="lastName"
					class="grow input-md"
					bind:value={payload.person.last_name}
					required
				/>
			</label>

			<label for="born" class="input input-bordered flex items-center gap-2">
				Born:
				<input
					type="date"
					id="born"
					class="grow input-md"
					bind:value={payload.person.born}
					required
				/>
			</label>

			<label for="birthplace" class="input input-bordered flex items-center gap-2">
				Birthplace:
				<input
					type="text"
					id="birthplace"
					class="grow input-md"
					bind:value={payload.person.birthplace}
					required
				/>
			</label>

			<label for="mothersFirstName" class="input input-bordered flex items-center gap-2">
				Mother's First Name:
				<input
					type="text"
					id="mothersFirstName"
					class="grow input-md"
					bind:value={payload.person.mothers_first_name}
					required
				/>
			</label>

			<label for="mothersLastName" class="input input-bordered flex items-center gap-2">
				Mother's Last Name:
				<input
					type="text"
					id="mothersLastName"
					class="grow input-md"
					bind:value={payload.person.mothers_last_name}
					required
				/>
			</label>

			<label for="death" class="input input-bordered flex items-center gap-2">
				Death:
				<input type="date" id="death" class="grow input-md" bind:value={payload.person.death} />
			</label>

			<label for="deathplace" class="input input-bordered flex items-center gap-2">
				Place of death:
				<input
					type="text"
					id="deathplace"
					class="grow input-md"
					bind:value={payload.person.deathplace}
				/>
			</label>

			<label for="residence" class="input input-bordered flex items-center gap-2">
				Residence:
				<input
					type="text"
					id="residence"
					class="grow input-md"
					bind:value={payload.person.residence}
				/>
			</label>

			<label for="titles" class="input input-bordered flex items-center gap-2">
				Titles:
				<input type="text" id="titles" class="grow input-md" bind:value={payload.person.titles} />
			</label>

			<button type="submit" class="btn btn-primary">Add</button>
		</form>
	</div>
</dialog>
