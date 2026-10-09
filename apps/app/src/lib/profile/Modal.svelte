<script lang="ts">
	import { died } from './../paraglide/messages/en.js';
	import { fade } from 'svelte/transition';

	import ModalButtons from './ModalButtons.svelte';
	import ProfileHeader from './ProfileHeader.svelte';
	import MediaGallery from './MediaGallery.svelte';
	import LifeEventsTimeline from './LifeEventsTimeline.svelte';
	import OtherDetails from './OtherDetails.svelte';
	import type { components } from '$lib/api/api.gen.js';
	import { deleteUnreferencedMedia } from './uploadMedia';
	import { savePerson } from './savePerson';

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
	const removedMedia = new Set<string>();

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
		removedMedia.clear();
	}

	function toggleEdit() {
		editorMode = !editorMode;
	}

	async function save() {
		try {
			console.debug('Saving person data:', draftPerson);
			const failure = await savePerson(person.id, draftPerson);
			if (failure) {
				console.error(failure);
				alert(failure);
				return;
			}

			person = { ...person, ...draftPerson };
			await deleteUnreferencedMedia(person.id!, removedMedia, person);
			removedMedia.clear();
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
		<ProfileHeader
			{person}
			{editorMode}
			onChange={handleDraftPersonChange}
			onRemoveMedia={(url) => removedMedia.add(url)}
		/>
		<MediaGallery
			{person}
			{editorMode}
			onChange={handleDraftPersonChange}
			onRemoveMedia={(url) => removedMedia.add(url)}
		/>
		<LifeEventsTimeline
			person_life_events={person.life_events}
			{editorMode}
			onChange={handleDraftPersonChange}
		/>
		<OtherDetails {person} {editorMode} onChange={handleDraftPersonChange} />
	</div>
</div>
