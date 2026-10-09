<script lang="ts">
	import type { components } from '$lib/api/api.gen';
	import { audio, video, photos, remove, upload } from '$lib/paraglide/messages';
	import UploadMediaModal from '$lib/profile/editors/UploadMediaModal.svelte';

	export let person: components['schemas']['PersonProperties'] & { id?: string };
	export let editorMode = false;
	export let onChange: (
		field: keyof components['schemas']['PersonProperties'],
		value: unknown
	) => void = () => {};
	export let onRemoveMedia: (url: string) => void = () => {};
	let uploadModal = false;
	let mediaType: 'audio' | 'video' | 'photo' = 'photo';

	function removeItem(field: 'photos' | 'videos' | 'audios', index: number) {
		const items = person[field] ?? [];
		const removed = items[index];
		person[field] = items.filter((_, i) => i !== index);
		onChange(field, person[field]);
		if (removed?.url) onRemoveMedia(removed.url);
	}
</script>

{#if person.photos?.length || person.videos?.length || person.audios?.length}
	<div class="divider">{photos()}, {video()} & {audio()}</div>
	<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
		{#each person.photos ?? [] as picture, index}
			<div class="relative">
				<img
					src={picture.url}
					alt={picture.description ?? photos()}
					class="h-32 w-full rounded-lg object-cover shadow-md"
				/>
				{#if editorMode}
					<button
						class="btn btn-circle btn-error btn-xs absolute top-1 right-1"
						aria-label={remove()}
						on:click={() => removeItem('photos', index)}>✕</button
					>
				{/if}
			</div>
		{/each}
		{#each person.videos ?? [] as video, index}
			<div class="relative">
				<video src={video.url} controls class="h-32 w-full rounded-lg shadow-md">
					<track kind="captions" />
				</video>
				{#if editorMode}
					<button
						class="btn btn-circle btn-error btn-xs absolute top-1 right-1"
						aria-label={remove()}
						on:click={() => removeItem('videos', index)}>✕</button
					>
				{/if}
			</div>
		{/each}
		{#each person.audios ?? [] as recording, index}
			<div class="relative">
				<audio src={recording.url} controls class="w-full">
					<track kind="captions" />
				</audio>
				{#if editorMode}
					<button
						class="btn btn-circle btn-error btn-xs absolute top-1 right-1"
						aria-label={remove()}
						on:click={() => removeItem('audios', index)}>✕</button
					>
				{/if}
			</div>
		{/each}
	</div>
{/if}

{#if editorMode}
	<div class="divider">{upload()}</div>
	<div class="grid grid-cols-3 gap-4">
		<button
			class="btn btn-soft btn-xs"
			on:click={() => {
				uploadModal = true;
				mediaType = 'photo';
			}}
		>
			{'+ ' + photos()}
		</button>
		<button
			class="btn btn-soft btn-xs"
			on:click={() => {
				uploadModal = true;
				mediaType = 'video';
			}}
		>
			{'+ ' + video()}
		</button>
		<button
			class="btn btn-soft btn-xs"
			on:click={() => {
				uploadModal = true;
				mediaType = 'audio';
			}}
		>
			{'+ ' + audio()}
		</button>
	</div>
{/if}

{#if uploadModal && person.id !== undefined}
	<UploadMediaModal
		closeModal={() => {
			uploadModal = false;
		}}
		{mediaType}
		personId={person.id}
		onCreation={(newMedia) => {
			if (mediaType === 'photo') {
				person.photos = [...(person.photos ?? []), newMedia];
				onChange('photos', person.photos);
			} else if (mediaType === 'video') {
				person.videos = [...(person.videos ?? []), newMedia];
				onChange('videos', person.videos);
			} else {
				person.audios = [...(person.audios ?? []), newMedia];
				onChange('audios', person.audios);
			}
		}}
	/>
{/if}
