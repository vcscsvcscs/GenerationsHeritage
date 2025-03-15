<script lang="ts">
	import '../app.css';
	import { i18n } from '$lib/i18n';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	let { children } = $props();
	import ThemeButton from '$lib/theme-select.svelte';
</script>

<ParaglideJS {i18n}>
	{@render children()}
	<div class="absolute top-2 right-2">
		<ThemeButton />
	</div>
</ParaglideJS>
