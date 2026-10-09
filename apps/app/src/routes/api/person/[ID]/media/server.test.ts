import { beforeEach, describe, expect, it, vi } from 'vitest';
import { DELETE, POST } from './+server';

const { GET } = vi.hoisted(() => ({ GET: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client: { GET } }));

const put = vi.fn();
const remove = vi.fn();
const KEY = 'people/1/123e4567-e89b-12d3-a456-426614174000.png';

function makeEvent(
	opts: {
		id?: string;
		session?: { userId: number } | null;
		bucket?: unknown;
		file?: File | null;
		kind?: string;
	} = {}
) {
	const body = new FormData();
	if (opts.file !== null) {
		body.append('file', opts.file ?? new File(['data'], 'photo.JPG', { type: 'image/jpeg' }));
	}
	return {
		params: { ID: opts.id ?? '1' },
		locals: { session: opts.session === undefined ? { userId: 1 } : opts.session },
		platform: { env: { GH_MEDIA: 'bucket' in opts ? opts.bucket : { put, delete: remove } } },
		url: new URL(`http://localhost/api/person/1/media${opts.kind ? `?kind=${opts.kind}` : ''}`),
		request: new Request('http://localhost/api/person/1/media', { method: 'POST', body })
	} as unknown as Parameters<typeof POST>[0];
}

function makeDeleteEvent(
	opts: { id?: string; key?: string; session?: { userId: number } | null; bucket?: unknown } = {}
) {
	return {
		params: { ID: opts.id ?? '1' },
		locals: { session: opts.session === undefined ? { userId: 1 } : opts.session },
		platform: { env: { GH_MEDIA: 'bucket' in opts ? opts.bucket : { put, delete: remove } } },
		url: new URL(`http://localhost/api/person/1/media?key=${opts.key ?? KEY}`)
	} as unknown as Parameters<typeof DELETE>[0];
}

describe('POST /api/person/[ID]/media', () => {
	beforeEach(() => {
		put.mockReset().mockResolvedValue({});
		remove.mockReset().mockResolvedValue(undefined);
		GET.mockReset().mockResolvedValue({ response: { ok: true, status: 200 } });
	});

	it('redirects to login without a session', async () => {
		await expect(POST(makeEvent({ session: null }))).rejects.toMatchObject({ status: 302 });
		expect(put).not.toHaveBeenCalled();
	});

	it('fails when the R2 binding is missing', async () => {
		const response = await POST(makeEvent({ bucket: undefined }));
		expect(response.status).toBe(500);
	});

	it('stores the file under people/{id}/{uuid}.{ext} for the profile owner', async () => {
		const response = await POST(makeEvent());

		expect(response.status).toBe(201);
		const { url, key } = (await response.json()) as { url: string; key: string };
		expect(key).toMatch(/^people\/1\/[0-9a-f-]{36}\.jpg$/);
		expect(url).toBe(`/api/media/${key}`);
		expect(put).toHaveBeenCalledWith(key, expect.any(File), {
			httpMetadata: { contentType: 'image/jpeg' }
		});
		expect(GET).not.toHaveBeenCalled();
	});

	it('checks admin rights when uploading for another person', async () => {
		const response = await POST(makeEvent({ id: '7' }));

		expect(response.status).toBe(201);
		expect(GET).toHaveBeenCalledWith('/admin/{id1}/{id2}', {
			params: { path: { id1: 7, id2: 1 }, header: { 'X-User-ID': 1 } }
		});
	});

	it.each([
		[401, 403],
		[403, 403],
		[404, 404],
		[500, 502],
		[503, 502]
	])('maps an /admin %i response to %i', async (adminStatus, expected) => {
		GET.mockResolvedValue({ response: { ok: false, status: adminStatus } });

		const response = await POST(makeEvent({ id: '7' }));

		expect(response.status).toBe(expected);
		expect(put).not.toHaveBeenCalled();
	});

	it('maps a network failure of the /admin call to 502', async () => {
		vi.spyOn(console, 'error').mockImplementation(() => {});
		GET.mockRejectedValue(new Error('network down'));

		const response = await POST(makeEvent({ id: '7' }));

		expect(response.status).toBe(502);
		expect(put).not.toHaveBeenCalled();
	});

	it('accepts images for kind=profile_picture', async () => {
		const response = await POST(makeEvent({ kind: 'profile_picture' }));
		expect(response.status).toBe(201);
	});

	it.each(['video/mp4', 'audio/mpeg'])('rejects %s for kind=profile_picture', async (type) => {
		const response = await POST(
			makeEvent({ kind: 'profile_picture', file: new File(['data'], 'f.bin', { type }) })
		);

		expect(response.status).toBe(415);
		expect(put).not.toHaveBeenCalled();
	});

	it('rejects an unknown kind', async () => {
		const response = await POST(makeEvent({ kind: 'avatar' }));
		expect(response.status).toBe(400);
	});

	it('rejects an invalid person id', async () => {
		const response = await POST(makeEvent({ id: 'abc' }));
		expect(response.status).toBe(400);
	});

	it('rejects a missing file', async () => {
		const response = await POST(makeEvent({ file: null }));
		expect(response.status).toBe(400);
	});

	it.each(['application/pdf', 'text/html', 'image/svg+xml'])('rejects %s', async (type) => {
		const response = await POST(makeEvent({ file: new File(['data'], 'f.bin', { type }) }));

		expect(response.status).toBe(415);
		expect(put).not.toHaveBeenCalled();
	});

	it('rejects files above the size limit', async () => {
		const big = new File([new Uint8Array(25 * 1024 * 1024 + 1)], 'big.mp4', { type: 'video/mp4' });

		const response = await POST(makeEvent({ file: big }));

		expect(response.status).toBe(413);
		expect(put).not.toHaveBeenCalled();
	});

	it('returns 500 when R2 fails', async () => {
		vi.spyOn(console, 'error').mockImplementation(() => {});
		put.mockRejectedValue(new Error('boom'));

		const response = await POST(makeEvent());

		expect(response.status).toBe(500);
	});
});

describe('DELETE /api/person/[ID]/media', () => {
	beforeEach(() => {
		remove.mockReset().mockResolvedValue(undefined);
		GET.mockReset().mockResolvedValue({ response: { ok: true, status: 200 } });
	});

	it('redirects to login without a session', async () => {
		await expect(DELETE(makeDeleteEvent({ session: null }))).rejects.toMatchObject({
			status: 302
		});
	});

	it('fails when the R2 binding is missing', async () => {
		expect((await DELETE(makeDeleteEvent({ bucket: undefined }))).status).toBe(500);
	});

	it('deletes an object under the person prefix', async () => {
		const response = await DELETE(makeDeleteEvent());

		expect(response.status).toBe(204);
		expect(remove).toHaveBeenCalledWith(KEY);
	});

	it('checks admin rights for another person and maps failures', async () => {
		const key = KEY.replace('people/1/', 'people/7/');
		GET.mockResolvedValue({ response: { ok: false, status: 401 } });

		const response = await DELETE(makeDeleteEvent({ id: '7', key }));

		expect(response.status).toBe(403);
		expect(remove).not.toHaveBeenCalled();
	});

	it.each([
		'people/2/123e4567-e89b-12d3-a456-426614174000.png',
		'people/1/../2/123e4567-e89b-12d3-a456-426614174000.png',
		'other/1/123e4567-e89b-12d3-a456-426614174000.png',
		'people/1/x.png',
		''
	])('rejects the key "%s" outside people/{ID}/', async (key) => {
		const response = await DELETE(makeDeleteEvent({ key }));

		expect(response.status).toBe(400);
		expect(remove).not.toHaveBeenCalled();
	});

	it('returns 500 when R2 fails', async () => {
		vi.spyOn(console, 'error').mockImplementation(() => {});
		remove.mockRejectedValue(new Error('boom'));

		expect((await DELETE(makeDeleteEvent())).status).toBe(500);
	});
});
