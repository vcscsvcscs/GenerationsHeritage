<script lang="ts">
    import * as msg from '$lib/paraglide/messages';
    export let draftPerson: any;
    export let editorMode = false;

    const skipFields = [
        'first_name',
        'last_name',
        'born',
        'biological_sex',
        'email',
        'limit',
        'mothers_first_name',
        'mothers_last_name',
        'profile_picture',
        'photos',
        'videos',
        'life_events'
    ];
</script>

<div class="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2">
    {#each Object.entries(draftPerson) as [key, value]}
        {#if !skipFields.includes(key)}
            <div>
                <p>
                    <strong>{(msg as any)[key]() ?? key}:</strong>
                    {#if editorMode}
                        <textarea
                            bind:value={draftPerson[key]}
                            class="textarea textarea-bordered textarea-sm w-full"
                        ></textarea>
                    {:else}
                        {JSON.stringify(value) ?? '-'}
                    {/if}
                </p>
            </div>
        {/if}
    {/each}
</div>
