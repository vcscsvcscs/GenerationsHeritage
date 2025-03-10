<script lang="ts">
	import { themes } from './themes';
	import { theme, light, dark, coffee } from '$lib/paraglide/messages.js';

	let current_theme = $state('');

	const themeMessages = new Map<string, string>([
		['light', light()],
		['dark', dark()],
		['coffee', coffee()],
		['cyberpunk', 'Cyberpunk'],
		['synthwave', 'Synthwave'],
		['retro', 'Retro'],
		['dracula', 'Dracula']
	]);

	$effect(() => {
		if (typeof window !== 'undefined') {
			const theme = window.localStorage.getItem('theme');
			if (theme && themes.includes(theme)) {
				document.documentElement.setAttribute('data-theme', theme);
				current_theme = theme;
			}
		}
	});

	function set_theme(event: Event) {
		const select = event.target as HTMLSelectElement;
		const theme = select.value;
		if (themes.includes(theme)) {
			const one_year = 60 * 60 * 24 * 365;
			window.localStorage.setItem('theme', theme);
			document.cookie = `theme=${theme}; max-age=${one_year}; path=/; SameSite=Lax`;
			document.documentElement.setAttribute('data-theme', theme);
			current_theme = theme;
		}
	}
</script>

<div class="dropdown mb-8">
	<select
		bind:value={current_theme}
		data-choose-theme
		class="select"
		onchange={set_theme}
	>
		<option value="" disabled={current_theme !== ''}>
			{theme()}
		</option>
		{#each themes as theme}
			<option
				value={theme}
				class="theme-controller capitalize"
				>{themeMessages.get(theme)}</option
			>
		{/each}
	</select>
</div>
