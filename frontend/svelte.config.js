import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: undefined,
			precompress: false,
			strict: true
		}),
		prerender: {
			handleHttpError: ({ path, referrer, message }) => {
				return {
					status: 404,
					html: `<h1>${message}</h1><p>Path: ${path}</p><p>Referrer: ${referrer}</p>`
				};
			}
		}
	}
};

export default config;
