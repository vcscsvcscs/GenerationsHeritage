<script lang="ts">
	import type { components } from "$lib/api/api.gen";
	import { video, photos } from "$lib/paraglide/messages";

    export let draftPerson: components['schemas']['PersonProperties'];
  </script>
  
  {#if draftPerson.photos?.length || draftPerson.videos?.length}
    <div class="divider">{photos()} & {video()}</div>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      {#each draftPerson.photos ?? [] as picture}
        <img src={picture.url} alt={picture.description ?? photos()} class="rounded-lg shadow-md object-cover w-full h-32" />
      {/each}
      {#each draftPerson.videos ?? [] as video}
        <video src={video.url} controls class="rounded-lg shadow-md w-full h-32">
            <track kind="captions" src={video.description} srcLang="en" default />
            <track kind="descriptions" src={video.description} srcLang="en" default />
        </video>
      {/each}
    </div>
  {/if}
  