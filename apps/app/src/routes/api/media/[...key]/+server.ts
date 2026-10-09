import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import { MEDIA_KEY_PATTERN } from '$lib/server/media';
import type { RequestEvent } from './$types';

function byteRange(range: R2Range, size: number): [number, number] {
	if ('suffix' in range) {
		return [Math.max(size - range.suffix, 0), size - 1];
	}
	const start = range.offset ?? 0;
	return [start, range.length === undefined ? size - 1 : start + range.length - 1];
}

export async function GET(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const bucket = event.platform?.env?.GH_MEDIA;
	if (!bucket) {
		return new Response('Server configuration error. GH_MEDIA R2 bucket missing', {
			status: 500
		});
	}

	const match = MEDIA_KEY_PATTERN.exec(event.params.key);
	if (!match) {
		return new Response(null, { status: 404 });
	}

	const personId = Number(match[1]);
	const userId = event.locals.session.userId;
	if (personId !== userId) {
		const person = await client.GET('/person/{id}', {
			params: {
				path: { id: personId },
				header: { 'X-User-ID': userId }
			}
		});
		if (!person.response.ok) {
			return new Response(null, { status: 403 });
		}
	}

	const object = await bucket.get(event.params.key, { range: event.request.headers });
	if (object === null) {
		return new Response(null, { status: 404 });
	}

	const headers = new Headers({
		'Cache-Control': 'private, max-age=31536000, immutable',
		'Accept-Ranges': 'bytes',
		'X-Content-Type-Options': 'nosniff',
		'Content-Security-Policy': "default-src 'none'; sandbox",
		ETag: object.httpEtag
	});
	object.writeHttpMetadata(headers);

	if (object.range) {
		const [start, end] = byteRange(object.range, object.size);
		headers.set('Content-Range', `bytes ${start}-${end}/${object.size}`);
		headers.set('Content-Length', String(end - start + 1));
		return new Response(object.body, { status: 206, headers });
	}

	headers.set('Content-Length', String(object.size));
	return new Response(object.body, { headers });
}
