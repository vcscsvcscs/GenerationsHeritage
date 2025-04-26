<script lang="ts">
	import { fade } from 'svelte/transition';

	import ModalButtons from './ModalButtons.svelte';
	import ProfileHeader from './ProfileHeader.svelte';
	import MediaGallery from './MediaGallery.svelte';
	import LifeEventsTimeline from './LifeEventsTimeline.svelte';
	import OtherDetails from './OtherDetails.svelte';
	import type { components } from '$lib/api/api.gen.js';
	let {
		open = false,
		person = {}
	}: { open: boolean; person: components['schemas']['PersonProperties'] } = $props();

	let editorMode = $state(false);
	let draftPerson = $state({});

	$effect(() => {
		if (open) {
			draftPerson = structuredClone(person);
			editorMode = false;
		}
	});

	function close() {
		open = false;
        editorMode = false;
        draftPerson = {};
	}

	function toggleEdit() {
		editorMode = !editorMode;
	}
</script>

{#if open}
	<div class="modal modal-open" transition:fade>
		<div class="modal-box max-h-screen w-full max-w-5xl overflow-y-auto">
			<div class="bg-base-100 sticky top-0 z-10">
				<ModalButtons {editorMode} onClose={close} onToggleEdit={toggleEdit} />
				<div class="divider"></div>
			</div>

			<ProfileHeader {draftPerson} {editorMode} />
			<MediaGallery {draftPerson} />
			<LifeEventsTimeline {draftPerson} />
			<OtherDetails {draftPerson} {editorMode} />
		</div>
	</div>
{/if}
