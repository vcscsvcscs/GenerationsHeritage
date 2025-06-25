<script lang="ts">
	import { died } from './../paraglide/messages/en.js';
	import { fade } from 'svelte/transition';

	import ModalButtons from './ModalButtons.svelte';
	import ProfileHeader from './ProfileHeader.svelte';
	import MediaGallery from './MediaGallery.svelte';
	import LifeEventsTimeline from './LifeEventsTimeline.svelte';
	import OtherDetails from './OtherDetails.svelte';
	import type { components } from '$lib/api/api.gen.js';

	let {
		closeModal = () => {},
		person = {}
	}: {
		closeModal: () => void;
		person: components['schemas']['PersonProperties'] & {
			id?: string;
		};
	} = $props();

	let editorMode = $state(false);
	let draftPerson = $state({} as components['schemas']['PersonProperties']);

	editorMode = false;

	function handleDraftPersonChange(
		field: keyof components['schemas']['PersonProperties'],
		value: any
	) {
		draftPerson[field] = value;
		if (field === 'invite_code') {
			save().then(() => {
				editorMode = true;
			});
			return;
		}
	}

	function close() {
		closeModal();
		editorMode = false;
		draftPerson = {};
	}

	function toggleEdit() {
		editorMode = !editorMode;
	}

	async function save() {
		try {
			console.debug('Saving person data:', draftPerson);
			const response = await fetch(`/api/person/${person.id}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(draftPerson)
			});

			if (!response.ok) {
				console.error('Error saving person data, status: ', response.status, (await response.json()));
				alert('Error saving person data, status: ' + response.status + (await response.json()));
				return;
			}

			if (response.status === 200) {
				person = { ...person, ...draftPerson };
				const data = (await response.json()) as {
					person?: components['schemas']['Person'];
				};
			} else {
				const errorDetails = await response.json();
				console.error('Error details:', errorDetails);
				alert(`Error saving person data, status: ${response.status} ${JSON.stringify(errorDetails)}`);
			}
		} catch (error) {
			alert('An unexpected error occurred: ' + error);
		}
		editorMode = !editorMode;
	}
</script>

<div class="modal modal-open" transition:fade>
	<div class="modal-box max-h-80 max-h-screen w-full max-w-5xl overflow-y-auto">
		<div class="bg-base-100 z-7 sticky top-0">
			<ModalButtons {editorMode} onClose={close} onSave={save} onToggleEdit={toggleEdit} />
			<div class="divider"></div>
		</div>
		<ProfileHeader {person} {editorMode} onChange={handleDraftPersonChange} />
		<MediaGallery {person} {editorMode} />
		<LifeEventsTimeline
			person_life_events={person.life_events}
			{editorMode}
			onChange={handleDraftPersonChange}
		/>
		<OtherDetails {person} {editorMode} onChange={handleDraftPersonChange} />
	</div>
</div>
