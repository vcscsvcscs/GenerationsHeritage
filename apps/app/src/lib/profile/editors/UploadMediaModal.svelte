<script lang="ts">
	import {
		audio,
		cancel,
		date,
		description,
		file,
		loading,
		media_title,
		missing_field,
		photos,
		upload,
		upload_failed,
		video
	} from '$lib/paraglide/messages';
	import { uploadMedia } from '$lib/profile/uploadMedia';

	export let closeModal: () => void;
	export let onCreation: (newMedia: {
		url: string;
		name: string;
		description: string;
		date?: string;
	}) => void = () => {};
	export let mediaType: 'audio' | 'video' | 'photo' = 'photo';
	export let personId: string | number;

	let selectedFile: File | null = null;
	let uploading = false;
	let error = '';

	let newMedia = {
		name: '',
		description: '',
		date: ''
	};

	$: acceptTypes =
		mediaType === 'audio' ? 'audio/*' : mediaType === 'video' ? 'video/*' : 'image/*';
	$: mediaLabel = mediaType === 'audio' ? audio() : mediaType === 'video' ? video() : photos();

	function handleFileChange(e: Event) {
		const input = e.target as HTMLInputElement;
		selectedFile = input.files?.[0] ?? null;
	}

	async function submit() {
		const name = newMedia.name.trim();
		if (!selectedFile) {
			error = missing_field({ field: file() });
			return;
		}
		if (!name) {
			error = missing_field({ field: media_title() });
			return;
		}

		uploading = true;
		error = '';
		try {
			const url = await uploadMedia(personId, selectedFile);
			onCreation({
				url,
				name,
				description: newMedia.description,
				...(newMedia.date && { date: newMedia.date })
			});
			closeModal();
		} catch (e) {
			error = e instanceof Error ? e.message : upload_failed();
		} finally {
			uploading = false;
		}
	}
</script>

<div class="modal modal-open z-8">
	<div class="modal-box w-full max-w-xl">
		<h3 class="text-lg font-bold">{upload()} {mediaLabel}</h3>

		<div class="form-control mt-4">
			<label for="mfile" class="label">{file()}</label>
			<input
				id="mfile"
				type="file"
				accept={acceptTypes}
				class="file-input file-input-bordered w-full"
				disabled={uploading}
				on:change={handleFileChange}
			/>
		</div>

		<div class="form-control mt-4">
			<label for="mtitle" class="label">{media_title()}</label>
			<input
				id="mtitle"
				bind:value={newMedia.name}
				class="input input-bordered w-full"
				disabled={uploading}
			/>
		</div>

		<div class="form-control mt-4">
			<label for="mdesc" class="label">{description()}</label>
			<textarea
				id="mdesc"
				bind:value={newMedia.description}
				class="textarea textarea-bordered w-full"
				disabled={uploading}
			>
			</textarea>
		</div>

		<div class="form-control mt-4">
			<label for="mdate" class="label">{date()}</label>
			<input
				id="mdate"
				type="date"
				bind:value={newMedia.date}
				class="input input-bordered w-full"
				disabled={uploading}
			/>
		</div>

		{#if error}
			<div role="alert" class="alert alert-error mt-4">{error}</div>
		{/if}

		<div class="modal-action">
			<button class="btn btn-outline" disabled={uploading} on:click={closeModal}>{cancel()}</button>
			<button class="btn btn-primary" disabled={uploading} on:click={submit}>
				{#if uploading}
					<span class="loading loading-spinner loading-xs"></span>
					{loading()}
				{:else}
					{upload()}
				{/if}
			</button>
		</div>
	</div>
</div>
