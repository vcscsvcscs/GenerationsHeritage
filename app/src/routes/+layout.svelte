<script lang="ts">
	import '../app.css';
	import { i18n } from '$lib/i18n';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import {title, family_tree} from '$lib/paraglide/messages.js';
	let { children } = $props();
</script>

<svelte:head>
	<title>{title({page:family_tree()})}</title>
</svelte:head>

<ParaglideJS {i18n}>
	{@render children()}
</ParaglideJS>
