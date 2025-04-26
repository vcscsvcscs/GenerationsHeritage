<script lang="ts">
    import type { components } from '$lib/api/api.gen';
    import { biological_sex, born, email, first_name, last_name, middle_name, mothers_first_name, mothers_last_name, profile_picture } from '$lib/paraglide/messages';
    export let draftPerson: components['schemas']['PersonProperties'];
    export let editorMode = false;
</script>
  
<div class="flex flex-col md:flex-row gap-6">
    {#if draftPerson.profile_picture}
      <div class="flex-shrink-0">
        <img src={draftPerson.profile_picture} alt={profile_picture()} class="rounded-lg shadow-md w-48 h-48 object-cover" />
      </div>
    {/if}
  
    <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <p><strong>{first_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.first_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.first_name ?? '-'}{/if}</p>
        <p><strong>{last_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.last_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.last_name ?? '-'}{/if}</p>
        <p><strong>{middle_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.middle_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.middle_name ?? '-'}{/if}</p>
        <p><strong>{born()}:</strong> {#if editorMode}<input bind:value={draftPerson.born} class="input input-sm input-bordered w-full" />{:else}{draftPerson.born ?? '-'}{/if}</p>
        <p><strong>{biological_sex()}:</strong> {#if editorMode}<input bind:value={draftPerson.biological_sex} class="input input-sm input-bordered w-full" />{:else}{draftPerson.biological_sex ?? '-'}{/if}</p>
      </div>
      <div>
        <p><strong>{email()}:</strong> {#if editorMode}<input bind:value={draftPerson.email} class="input input-sm input-bordered w-full" />{:else}{draftPerson.email ?? '-'}{/if}</p>
        <p><strong>Limit:</strong> {#if editorMode}<input bind:value={draftPerson.limit} class="input input-sm input-bordered w-full" />{:else}{draftPerson.limit ?? '-'}{/if}</p>
        <p><strong>{mothers_first_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.mothers_first_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.mothers_first_name ?? '-'}{/if}</p>
        <p><strong>{mothers_last_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.mothers_last_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.mothers_last_name ?? '-'}{/if}</p>
      </div>
    </div>
</div>
  