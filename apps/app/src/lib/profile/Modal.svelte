<script lang="ts">
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
	}: { closeModal: () => void; person: components['schemas']['PersonProperties'] } = $props();

	let editorMode = $state(false);
	let draftPerson = $state({} as components['schemas']['PersonProperties']);

	editorMode = false;

	function handleDraftPersonChange(
		field: keyof components['schemas']['PersonProperties'],
		value: any
	) {
		draftPerson[field] = value;
	}

	function close() {
		closeModal();
		editorMode = false;
		draftPerson = {};
	}

	function toggleEdit() {
		editorMode = !editorMode;
	}

	function save() {
		// Save logic here
		editorMode = false;
	}
</script>

<div class="modal modal-open" transition:fade>
	<div class="modal-box max-h-screen w-full max-w-5xl overflow-y-auto">
		<div class="bg-base-100 sticky top-0 z-10">
			<ModalButtons {editorMode} onClose={close} onSave={save} onToggleEdit={toggleEdit} />
			<div class="divider"></div>
		</div>
		<ProfileHeader {person} {editorMode} onChange={handleDraftPersonChange} />
		<MediaGallery {person} />
		<LifeEventsTimeline
			person_life_events={person.life_events}
			{editorMode}
			onChange={handleDraftPersonChange}
		/>
		<OtherDetails {person} {editorMode} onChange={handleDraftPersonChange} />
	</div>
</div>
