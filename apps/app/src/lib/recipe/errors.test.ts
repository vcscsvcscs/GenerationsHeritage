import { describe, it, expect } from 'vitest';
import { errorMessage } from './errors';

describe('errorMessage', () => {
	it('distinguishes unauthorized, not found and generic failures', () => {
		const messages = [401, 404, 500].map(errorMessage);
		expect(new Set(messages).size).toBe(3);
		expect(errorMessage(0)).toBe(errorMessage(500));
		expect(errorMessage(409)).toBe(errorMessage(500));
	});
});
