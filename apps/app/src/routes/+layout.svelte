<script lang="ts">
	import '../app.css';
	import { i18n } from '$lib/i18n';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	let { children } = $props();
	import ThemeButton from '$lib/theme-select.svelte';
	import Logout from '$lib/logout.svelte';
	import { page } from '$app/state';
</script>

<ParaglideJS {i18n}>
	{@render children()}
	<div class="flex flex-row absolute top-2 right-2">
		<ThemeButton />
		<Logout show={!page.url.pathname.includes("login")}/>
	</div>
</ParaglideJS>
