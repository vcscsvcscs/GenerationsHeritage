<script lang="ts">
	import { life_events, unknown, until } from '$lib/paraglide/messages';

	export let draftPerson: any;
</script>

{#if draftPerson.life_events?.length}
	<div class="divider">{life_events()}</div>
	<ul class="timeline timeline-snap-start timeline-vertical">
		{#each draftPerson.life_events as event}
			<li>
				<div class="timeline-start">{event.from ?? unknown()}</div>
				<div class="timeline-middle">
					<div class="badge badge-primary"></div>
				</div>
				<div class="timeline-end">
					<p>{event.description}</p>
					{#if event.to}
						<p class="text-sm opacity-50">{until()} {event.to}</p>
					{/if}
				</div>
				<hr />
			</li>
		{/each}
	</ul>
{/if}
