<script lang="ts">
	import { logout } from '$lib/paraglide/messages.js';
	export let show = false;
</script>

{#if show}
	<a class="btn btn-error btn-xs h-8 min-h-0 px-4 py-0 text-sm" href="/logout">{logout()}</a>
{/if}
