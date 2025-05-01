import { describe, it, expect, vi } from 'vitest';
import type { Mock } from 'vitest';
import { switchToLanguage } from './switchToLanguage';
import { i18n } from '$lib/i18n';
import { goto } from '$app/navigation';

vi.mock('$lib/i18n', () => ({
	i18n: {
		route: vi.fn().mockImplementation((translatedPath: string) => ''),
		resolveRoute: vi.fn().mockImplementation((path: string, lang?: string) => '')
	}
}));

vi.mock('$app/state', () => ({
	page: {
		url: {
			pathname: '/current-path'
		}
	}
}));

vi.mock('$app/navigation', () => ({
	goto: vi.fn()
}));

describe('switchToLanguage', () => {
	it('should switch to the new language', () => {
		const newLanguage = 'en';
		const canonicalPath = '/canonical-path';
		const localisedPath = '/en/canonical-path';

		(i18n.route as Mock).mockReturnValue(canonicalPath);
		(i18n.resolveRoute as Mock).mockReturnValue(localisedPath);

		switchToLanguage(newLanguage);

		expect(i18n.route).toHaveBeenCalledWith('/current-path');
		expect(i18n.resolveRoute).toHaveBeenCalledWith(canonicalPath, newLanguage);
		expect(goto).toHaveBeenCalledWith(localisedPath);
	});
});
