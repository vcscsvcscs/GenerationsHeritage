import { afterEach, describe, expect, it, vi } from 'vitest';
import { savePerson } from './savePerson';

describe('savePerson', () => {
	afterEach(() => {
		vi.unstubAllGlobals();
	});

	it('PATCHes the full draft and returns null on success', async () => {
		const fetchMock = vi.fn().mockResolvedValue(Response.json({}));
		vi.stubGlobal('fetch', fetchMock);
		const draft = { photos: [{ url: '/api/media/people/3/a.png' }] };

		await expect(savePerson('3', draft)).resolves.toBeNull();

		expect(fetchMock).toHaveBeenCalledWith('/api/person/3', {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(draft)
		});
	});

	it('reads the error body exactly once and reports status and details', async () => {
		const response = Response.json({ msg: 'nope' }, { status: 401 });
		const json = vi.spyOn(response, 'json');
		vi.stubGlobal('fetch', vi.fn().mockResolvedValue(response));

		const failure = await savePerson('3', {});

		expect(json).toHaveBeenCalledTimes(1);
		expect(failure).toBe('Error saving person data, status: 401 {"msg":"nope"}');
	});

	it('still reports the status when the error body is not JSON', async () => {
		vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('oops', { status: 500 })));

		await expect(savePerson('3', {})).resolves.toBe('Error saving person data, status: 500 null');
	});
});
