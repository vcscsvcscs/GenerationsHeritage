import { afterEach, describe, expect, it, vi } from 'vitest';
import { uploadMedia } from './uploadMedia';

const file = new File(['data'], 'me.png', { type: 'image/png' });

function mockFetch(response: Response | Error) {
	const fetchMock =
		response instanceof Error
			? vi.fn().mockRejectedValue(response)
			: vi.fn().mockResolvedValue(response);
	vi.stubGlobal('fetch', fetchMock);
	return fetchMock;
}

describe('uploadMedia', () => {
	afterEach(() => {
		vi.unstubAllGlobals();
		vi.restoreAllMocks();
	});

	it('posts the file as multipart form data and returns the url', async () => {
		const fetchMock = mockFetch(Response.json({ url: '/api/media/people/3/x.png' }));

		await expect(uploadMedia(3, file)).resolves.toBe('/api/media/people/3/x.png');

		const [path, init] = fetchMock.mock.calls[0];
		expect(path).toBe('/api/person/3/media');
		expect(init.method).toBe('POST');
		expect((init.body as FormData).get('file')).toBeInstanceOf(File);
	});

	it.each([413, 415, 500])('rejects with a message for status %i', async (status) => {
		mockFetch(new Response(null, { status }));

		await expect(uploadMedia(3, file)).rejects.toThrow(/.+/);
	});

	it('uses distinct messages for too large and unsupported files', async () => {
		const messages: string[] = [];
		for (const status of [413, 415, 500]) {
			mockFetch(new Response(null, { status }));
			messages.push(await uploadMedia(3, file).catch((e: Error) => e.message));
		}

		expect(new Set(messages).size).toBe(3);
	});

	it('rejects with the upload failed message on network errors', async () => {
		vi.spyOn(console, 'error').mockImplementation(() => {});
		mockFetch(new TypeError('Failed to fetch'));

		await expect(uploadMedia(3, file)).rejects.toThrow(/.+/);
	});
});
