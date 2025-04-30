<script lang="ts">
	import type { components } from '$lib/api/api.gen';
	import { video, photos, upload } from '$lib/paraglide/messages';
	import UploadMediaModal from '$lib/profile/editors/UploadMediaModal.svelte';

	export let person: components['schemas']['PersonProperties'];
	export let editorMode = false;
	let uploadModal = false;
	let mediaType: 'audio' | 'video' | 'photo' | undefined = undefined;
</script>

{#if person.photos?.length || person.videos?.length}
	<div class="divider">{photos()} & {video()}</div>
	<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
		{#each person.photos ?? [] as picture}
			<img
				src={picture.url}
				alt={picture.description ?? photos()}
				class="h-32 w-full rounded-lg object-cover shadow-md"
			/>
		{/each}
		{#each person.videos ?? [] as video}
			<video src={video.url} controls class="h-32 w-full rounded-lg shadow-md">
				<track kind="captions" src={video.description} srcLang="en" default />
				<track kind="descriptions" src={video.description} srcLang="en" default />
			</video>
		{/each}
	</div>
{/if}

{#if editorMode}
	<div class="divider">{upload()}</div>
	<div class="grid grid-cols-2 gap-4">
		<button
			class="btn btn-soft btn-xs"
			on:click={() => {
				uploadModal = true;
				mediaType = 'photo';
			}}
		>
			{'+ '+photos()}
		</button>
		<button
			class="btn btn-soft btn-xs"
			on:click={() => {
				uploadModal = true;
				mediaType = 'video';
			}}
		>
			{'+ '+video()}
		</button>
	</div>
{/if}

{#if uploadModal}
	<UploadMediaModal
		closeModal={() => {
			uploadModal = false;
		}}
		{mediaType}
		onCreation={(newMedia: { url: string; name: string; description: string; date: string }) => {
			if (mediaType === 'photo') {
				person.photos = [...(person.photos ?? []), newMedia];
			} else if (mediaType === 'video') {
				person.videos = [...(person.videos ?? []), newMedia];
			}
		}}
	/>
{/if}
