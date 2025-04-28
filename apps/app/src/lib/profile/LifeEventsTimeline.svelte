<script lang="ts">
	import {
		add_life_event,
		description,
		life_events,
		unknown,
		until
	} from '$lib/paraglide/messages';
	import type { components } from '$lib/api/api.gen';

	export let person_life_events: components['schemas']['PersonProperties']['life_events'];
	export let editorMode = false;
	export let onChange: (field: keyof components['schemas']['PersonProperties'], value: any) => void;

	function updateEvent(index: number, key: 'from' | 'to' | 'description', value: string) {
		if (!person_life_events) return;

		person_life_events = person_life_events.map((event, i) =>
			i === index ? { ...event, [key]: value } : event
		);

		onChange('life_events', person_life_events);
	}

	function addEvent() {
		const newEvent = { from: '', to: '', description: '' };
		person_life_events = [...(person_life_events ?? []), newEvent];
		onChange('life_events', person_life_events);
	}
</script>

{#if person_life_events?.length}
	<div class="divider">{life_events()}</div>
	<ul class="timeline timeline-snap-start timeline-vertical">
		{#each person_life_events as event, index}
			<li>
				<div class="timeline-start">
					{#if editorMode}
						<input
							type="text"
							class="input input-xs input-bordered"
							value={event.from ?? ''}
							on:input={(e) => updateEvent(index, 'from', e.currentTarget.value)}
							placeholder={unknown().toLowerCase()}
						/>
					{:else}
						{event.from ?? unknown().toLowerCase()}
					{/if}
				</div>

				<div class="timeline-middle">
					<div class="badge badge-primary"></div>
				</div>

				<div class="timeline-end space-y-1">
					{#if editorMode}
						<textarea
							class="textarea textarea-xs textarea-bordered w-full"
							value={event.description ?? ''}
							on:input={(e) => updateEvent(index, 'description', e.currentTarget.value)}
							placeholder={description()}
						></textarea>
					{:else}
						<p>{event.description}</p>
					{/if}

					{#if event.to || editorMode}
						<p class="text-sm opacity-50">
							{until()}
							{#if editorMode}
								<input
									type="text"
									class="input input-xs input-bordered ml-1"
									value={event.to ?? ''}
									on:input={(e) => updateEvent(index, 'to', e.currentTarget.value)}
									placeholder={unknown().toLowerCase()}
								/>
							{:else}
								{event.to ?? unknown().toLowerCase()}
							{/if}
						</p>
					{/if}
				</div>
				<hr />
			</li>
		{/each}
	</ul>
{/if}

{#if editorMode}
	<div class="mt-4 flex justify-center">
		<button class="btn btn-primary btn-sm" on:click={addEvent}>
			{add_life_event()}
		</button>
	</div>
{/if}
