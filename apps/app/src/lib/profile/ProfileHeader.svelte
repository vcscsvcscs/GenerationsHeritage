<script lang="ts">
	import { death } from './../paraglide/messages/en.js';
  import { onMount } from 'svelte';
    import type { components } from '$lib/api/api.gen';
    import { male,female,intersex,other,change_profile_picture, biological_sex, born,died, email, first_name, id, last_name, middle_name, mothers_first_name, mothers_last_name, profile_picture } from '$lib/paraglide/messages';
    import { callMessageFunction } from '$lib/i18n';
    import type { MessageKeys } from '$lib/i18n';

    export let draftPerson: components['schemas']['PersonProperties'] & {
      id?: string,
    };
    export let editorMode = false;
    let birth_date: HTMLInputElement;
    let death_date: HTMLInputElement;
    onMount(() => {
		if (birth_date) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: birth_date,
					onSelect: function (date) {
						birth_date.value = date.toISOString();
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
    if (death_date) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: death_date,
					onSelect: function (date) {
						death_date.value = date.toISOString();
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
	});
</script>
  
<div class="flex flex-col md:flex-row gap-6">
    <div class="flex-shrink-0 flex flex-col items-center gap-2">
      <img src={draftPerson.profile_picture||'https://cdn-icons-png.flaticon.com/512/10628/10628885.png'} alt={profile_picture()} class="rounded-lg shadow-md w-48 h-48 object-cover" />  
      {#if editorMode}
        <button class="btn bg-neutral text-neutral-content btn-xs" on:click={() => {}}>
          {change_profile_picture()}
        </button>
      {/if}
    </div>
    <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <p><strong>{first_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.first_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.first_name ?? '-'}{/if}</p>
        <p><strong>{last_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.last_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.last_name ?? '-'}{/if}</p>
        <p><strong>{middle_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.middle_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.middle_name ?? '-'}{/if}</p>
        <p><strong>{born()}:</strong>         
              <input
              type="text"
              class="w-full pika-single"
              id="birth_date"
              bind:this={birth_date}
              bind:value={draftPerson.born}/> 
        </p>
        <p><strong>{died()}:</strong>
            <input
            type="text"
            class="w-full pika-single"
            id="death_date"
            placeholder={died()}
            bind:this={death_date}
            bind:value={draftPerson.died}
          />
        </p>
        <p><strong>{biological_sex()}:</strong>
          {#if editorMode}
          <select
							name="biological_sex"
							class="select select-bordered w-full select-sm"
							id="biological_sex"
              bind:value={draftPerson.biological_sex}
							placeholder={biological_sex()}
						>
							<option value="male">{male()} </option>
							<option value="female">{female()} </option>
							<option value="intersex">{intersex()} </option>
							<option value="other">{other()} </option>
						</select>
          {:else}{callMessageFunction(draftPerson.biological_sex as MessageKeys) ?? '-'}{/if}</p>
      </div>
      <div>
        <p><strong>{email()}:</strong> {#if editorMode}<input bind:value={draftPerson.email} class="input input-sm input-bordered w-full" />{:else}{draftPerson.email ?? '-'}{/if}</p>
        <p><strong>{mothers_first_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.mothers_first_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.mothers_first_name ?? '-'}{/if}</p>
        <p><strong>{mothers_last_name()}:</strong> {#if editorMode}<input bind:value={draftPerson.mothers_last_name} class="input input-sm input-bordered w-full" />{:else}{draftPerson.mothers_last_name ?? '-'}{/if}</p>
        <p><strong> {id()}:</strong>{draftPerson.id ?? '-'}</p>
        <p><strong> Limit:</strong>{draftPerson.limit ?? '-'}</p>
      </div>
    </div>
</div>
  