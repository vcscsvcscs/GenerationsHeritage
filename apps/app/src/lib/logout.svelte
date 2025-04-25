<script lang="ts">
	import { logout } from '$lib/paraglide/messages.js';
	export let show = false;
</script>

{#if show}
	<a class="btn btn-error btn-xs" href="/logout">{logout()}</a>
{/if}