<script lang="ts">
	import { logout } from '$lib/paraglide/messages.js';
	export let show = false;
</script>

{#if show}
	<a class="btn btn-error btn-xs min-h-0 h-8 px-4 py-0 text-sm" href="/logout">{logout()}</a>
{/if}