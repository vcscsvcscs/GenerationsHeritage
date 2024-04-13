module.exports = {
	extends: ['plugin:svelte/base', 'plugin:svelte/recommended', 'plugin:svelte/prettier'],
	overrides: [
		{
			files: ['*.svelte'],
			parser: 'svelte-eslint-parser'
		}
	],
	env: {
		es6: true
	},
	parserOptions: {
		ecmaFeatures: {
			experimentalObjectRestSpread: true,
			jsx: true
		},
		sourceType: 'module'
	}
};
