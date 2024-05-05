<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { setFamilyTreeNodes } from './setFamilyTreeNodes';
	import { user } from '$lib/stores';
	export let showPanel = false;
	export let id = '';
	let relationship = '';
	let second_person_id = '';
	let dialog: HTMLDialogElement; // HTMLDialogElement
	let auth_token = '';

	user.subscribe((value) => {
		if (value) {
			auth_token = value.access_token;
		}
	});

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		if (second_person_id && id && relationship) {
			const payload = {
				first_person_id: id,
				second_person_id: second_person_id,
				relationship: relationship,
				direction: '->'
			};

			fetch(PUBLIC_API_URL + '/relationship', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: auth_token
				},
				body: JSON.stringify(payload)
			}).then((response) => {
				if (response.ok) {
					showPanel = false;
					setFamilyTreeNodes();
					dialog.close();
				} else {
					alert('Failed to add relationship! with status ' + response.status);
				}
			});
		} else {
			alert('You have to fill out all the fields!');
		}
	}

	$: if (dialog && showPanel) dialog.showModal();
</script>

<dialog bind:this={dialog} class="modal bg-base-300">
	<div class="modal-box w-11/12 max-w-5xl">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
		</form>
		<h1>Add a relationship</h1>
		<p>
			Ask your relative to give you their id and set what kind of relationship you have with them.
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

			<button type="submit" class="btn btn-primary">Add</button>
		</form>
	</div>
</dialog>
