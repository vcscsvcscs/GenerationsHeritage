import { paraglide } from '@inlang/paraglide-sveltekit/vite';
import { defineConfig } from 'vitest/config';
import { execSync } from "child_process";
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
	plugins: [
		{
			name: "openapi-generate",
			buildStart() {
			  console.log("Generating TypeScript client from OpenAPI...");
			  execSync("npx openapi-typescript ../../api/openapi.json --output src/lib/api/api.gen.ts", { stdio: "inherit" });
			  console.log("OpenAPI client generated!");
			}
		},
		sveltekit(),
		paraglide({
			project: './project.inlang',
			outdir: './src/lib/paraglide'
		})
	],

	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	}
});
