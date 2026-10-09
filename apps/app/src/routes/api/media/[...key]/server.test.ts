import { beforeEach, describe, expect, it, vi } from 'vitest';
import { GET as serve } from './+server';

const { GET } = vi.hoisted(() => ({ GET: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client: { GET } }));

const KEY = 'people/1/123e4567-e89b-12d3-a456-426614174000.mp4';
const get = vi.fn();

function makeEvent(opts: { key?: string; userId?: number; range?: string } = {}) {
	return {
		params: { key: opts.key ?? KEY },
		locals: { session: { userId: opts.userId ?? 1 } },
		platform: { env: { GH_MEDIA: { get } } },
		request: new Request('http://localhost/api/media/x', {
			headers: opts.range ? { Range: opts.range } : {}
		})
	} as unknown as Parameters<typeof serve>[0];
}

function r2Object(range?: R2Range) {
	return {
		body: 'body',
		size: 100,
		range,
		httpEtag: '"etag"',
		writeHttpMetadata: (headers: Headers) => headers.set('Content-Type', 'video/mp4')
	};
}

describe('GET /api/media/[...key]', () => {
	beforeEach(() => {
		get.mockReset().mockResolvedValue(r2Object());
		GET.mockReset().mockResolvedValue({ response: { ok: true } });
	});

	it('returns 404 for keys outside the people/{id}/{uuid}.{ext} scheme', async () => {
		const response = await serve(makeEvent({ key: 'people/1/../2/x.png' }));
		expect(response.status).toBe(404);
		expect(get).not.toHaveBeenCalled();
	});

	it('streams the object with content type and cache headers for the owner', async () => {
		const response = await serve(makeEvent());

		expect(response.status).toBe(200);
		expect(response.headers.get('Content-Type')).toBe('video/mp4');
		expect(response.headers.get('Cache-Control')).toContain('immutable');
		expect(response.headers.get('Content-Length')).toBe('100');
		expect(GET).not.toHaveBeenCalled();
	});

	it('checks the profile with the db-adapter for other users', async () => {
		await serve(makeEvent({ userId: 2 }));

		expect(GET).toHaveBeenCalledWith('/person/{id}', {
			params: { path: { id: 1 }, header: { 'X-User-ID': 2 } }
		});
	});

	it('returns 403 when the user can not see the profile', async () => {
		GET.mockResolvedValue({ response: { ok: false, status: 401 } });

		const response = await serve(makeEvent({ userId: 2 }));

		expect(response.status).toBe(403);
		expect(get).not.toHaveBeenCalled();
	});

	it('returns 404 when the object does not exist', async () => {
		get.mockResolvedValue(null);
		expect((await serve(makeEvent())).status).toBe(404);
	});

	it('answers range requests with 206 and Content-Range', async () => {
		get.mockResolvedValue(r2Object({ offset: 10, length: 20 }));

		const response = await serve(makeEvent({ range: 'bytes=10-29' }));

		expect(response.status).toBe(206);
		expect(response.headers.get('Content-Range')).toBe('bytes 10-29/100');
		expect(response.headers.get('Content-Length')).toBe('20');
		expect(get.mock.calls[0][1].range.get('Range')).toBe('bytes=10-29');
	});
});
