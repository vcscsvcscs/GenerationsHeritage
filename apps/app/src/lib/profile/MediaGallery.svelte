<script lang="ts">
	import type { components } from '$lib/api/api.gen';
	import { video, photos, upload } from '$lib/paraglide/messages';

	export let draftPerson: components['schemas']['PersonProperties'];
  export let editorMode = false;
</script>

{#if editorMode}
  <button class="btn bg-neutral text-neutral-content btn-xs" on:click={() => {}}>
    {upload()}
  </button>
{/if}

{#if draftPerson.photos?.length || draftPerson.videos?.length}
	<div class="divider">{photos()} & {video()}</div>
	<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
		{#each draftPerson.photos ?? [] as picture}
			<img
				src={picture.url}
				alt={picture.description ?? photos()}
				class="h-32 w-full rounded-lg object-cover shadow-md"
			/>
		{/each}
		{#each draftPerson.videos ?? [] as video}
			<video src={video.url} controls class="h-32 w-full rounded-lg shadow-md">
				<track kind="captions" src={video.description} srcLang="en" default />
				<track kind="descriptions" src={video.description} srcLang="en" default />
			</video>
		{/each}
	</div>
{/if}
