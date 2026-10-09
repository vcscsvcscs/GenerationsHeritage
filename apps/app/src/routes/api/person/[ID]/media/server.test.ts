import { beforeEach, describe, expect, it, vi } from 'vitest';
import { POST } from './+server';

const { GET } = vi.hoisted(() => ({ GET: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client: { GET } }));

const put = vi.fn();

function makeEvent(
	opts: {
		id?: string;
		session?: { userId: number } | null;
		bucket?: unknown;
		file?: File | null;
	} = {}
) {
	const body = new FormData();
	if (opts.file !== null) {
		body.append('file', opts.file ?? new File(['data'], 'photo.JPG', { type: 'image/jpeg' }));
	}
	return {
		params: { ID: opts.id ?? '1' },
		locals: { session: opts.session === undefined ? { userId: 1 } : opts.session },
		platform: { env: { GH_MEDIA: 'bucket' in opts ? opts.bucket : { put } } },
		request: new Request('http://localhost/api/person/1/media', { method: 'POST', body })
	} as unknown as Parameters<typeof POST>[0];
}

describe('POST /api/person/[ID]/media', () => {
	beforeEach(() => {
		put.mockReset().mockResolvedValue({});
		GET.mockReset().mockResolvedValue({ response: { ok: true } });
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

	it('rejects users that can not manage the person', async () => {
		GET.mockResolvedValue({ response: { ok: false, status: 500 } });

		const response = await POST(makeEvent({ id: '7' }));

		expect(response.status).toBe(403);
		expect(put).not.toHaveBeenCalled();
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
