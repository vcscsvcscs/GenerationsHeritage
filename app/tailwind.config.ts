import type { Config } from 'tailwindcss';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {}
	},
	daisyui: {
		themes: ['light', 'dark', 'cyberpunk', 'synthwave', 'retro', 'roboto', 'dracula']
	},
	plugins: [require('tailwindcss'), require('autoprefixer')]
} satisfies Config;
