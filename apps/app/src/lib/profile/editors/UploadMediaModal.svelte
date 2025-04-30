<script lang="ts">
	import { date, description, file, media_title, title, upload } from '$lib/paraglide/messages';

	export let closeModal: () => void;
	export let onCreation: (newMedia: {
		url: string;
		name: string;
		description: string;
		date: string;
	}) => void = () => {};
	export let mediaType: 'audio' | 'video' | 'photo' = 'photo';

	let selectedFile: File | null = null;

	let newMedia = {
		url: '',
		name: '',
		description: '',
		date: ''
	};

	// Determine accepted input types based on mediaType
	$: acceptTypes =
		mediaType === 'audio' ? 'audio/*' : mediaType === 'video' ? 'video/*' : 'image/*';

	function handleFileChange(e: Event) {
		const input = e.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			selectedFile = input.files[0];
		}
	}

	async function uploadMedia() {
		if (!selectedFile) {
			alert('Please select a file');
			return;
		}

		// Simulate file upload (replace with actual upload logic)
		newMedia.url = URL.createObjectURL(selectedFile);

		// Emit event using custom dispatch
		const uploadEvent = new CustomEvent('upload', {
			detail: { ...newMedia }
		});
		dispatchEvent(uploadEvent);

		// Clean up
		selectedFile = null;
		newMedia = { url: '', name: '', description: '', date: '' };
		onCreation(newMedia);
		closeModal();
	}
</script>

<div class="modal modal-open z-8">
	<div class="modal-box w-full max-w-xl">
		<h3 class="text-lg font-bold">{upload() + mediaType}</h3>

		<div class="form-control mt-4">
			<label for="mfile" class="label">{upload() + ' ' + file()}</label>
			<input
				id="mfile"
				type="file"
				accept={acceptTypes}
				class="file-input file-input-bordered w-full"
				on:change={handleFileChange}
			/>
		</div>

		<div class="form-control mt-4">
			<label for="mtitle" class="label">{media_title()}</label>
			<input id="mtitle" bind:value={newMedia.name} class="input input-bordered w-full" />
		</div>

		<div class="form-control mt-4">
			<label for="mdesc" class="label">{description()}</label>
			<textarea
				id="mdesc"
				bind:value={newMedia.description}
				class="textarea textarea-bordered w-full"
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
			/>
		</div>

		<div class="modal-action">
			<button class="btn btn-outline" on:click={closeModal}>Cancel</button>
			<button class="btn btn-primary" on:click={uploadMedia}>Upload</button>
		</div>
	</div>
</div>
