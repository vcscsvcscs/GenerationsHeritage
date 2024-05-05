<script>
	import { PUBLIC_API_URL } from '$env/static/public';
	import { setFamilyTreeNodes } from './setFamilyTreeNodes';
	export let showPanel = false;
	export let id = '';
	let relationship = '';
	let dialog; // HTMLDialogElement

	let firstName = '';
	let middleName = '';
	let lastName = '';
	let born = '';
	let mothersFirstName = '';
	let mothersLastName = '';

	function handleSubmit(event) {
		event.preventDefault();

		if (
			id &&
			relationship &&
			firstName &&
			lastName &&
			born &&
			mothersFirstName &&
			mothersLastName
		) {
			const payload = {
				relationship: {
					first_person_id: id,
					relationship: relationship,
					direction: '->'
				},
				person: {
					first_name: firstName,
					middle_name: middleName,
					last_name: lastName,
					mothers_first_name: mothersFirstName,
					mothers_last_name: mothersLastName,
					born: born
				}
			};

			fetch(PUBLIC_API_URL + '/createRelationshipAndPerson', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
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
			alert('You have to fill out all the required fields!');
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

			<label for="firstName">First Name:</label>
			<input type="text" id="firstName" bind:value={firstName} required />

			<label for="middleName">Middle Name:</label>
			<input type="text" id="middleName" bind:value={middleName} />

			<label for="lastName">Last Name:</label>
			<input type="text" id="lastName" bind:value={lastName} required />

			<label for="born">Born:</label>
			<input type="date" id="born" bind:value={born} required />

			<label for="mothersFirstName">Mother's First Name:</label>
			<input type="text" id="mothersFirstName" bind:value={mothersFirstName} required />

			<label for="mothersLastName">Mother's Last Name:</label>
			<input type="text" id="mothersLastName" bind:value={mothersLastName} required />

			<button type="submit" class="btn btn-primary">Add</button>
		</form>
	</div>
</dialog>
