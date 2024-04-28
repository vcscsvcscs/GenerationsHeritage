/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,svelte,js,ts}'],
	theme: {
		extend: {}
	},
	daisyui: {
		themes: ['light', 'dark', 'cyberpunk', 'synthwave', 'retro', 'roboto', 'dracula']
	},
	plugins: [require('daisyui')]
};
