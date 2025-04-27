import * as runtime from '$lib/paraglide/runtime';
import { createI18n } from '@inlang/paraglide-sveltekit';

export const i18n = createI18n(runtime);

import * as messages from '$lib/paraglide/messages';

export type MessageKeys = keyof typeof messages;

export function callMessageFunction(name: MessageKeys): string {
	const fn = messages[name];
	if (typeof fn === 'function') {
		return fn({ thing: '', field: '', page: '', name: '' });
	} else {
		throw new Error(`Function ${name} is not callable`);
	}
}
