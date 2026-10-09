import { redirect } from '@sveltejs/kit';
import { v4 as uuidv4 } from 'uuid';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';

const MAX_MEDIA_SIZE = 25 * 1024 * 1024;

const ALLOWED_TYPE = /^(image|video|audio)\/(?!svg)[\w.+-]+$/;

function json(body: unknown, status: number): Response {
	return new Response(JSON.stringify(body), { status });
}

export async function POST(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const bucket = event.platform?.env?.GH_MEDIA;
	if (!bucket) {
		return new Response('Server configuration error. GH_MEDIA R2 bucket missing', {
			status: 500
		});
	}

	const personId = Number(event.params.ID);
	const userId = event.locals.session.userId;
	if (!Number.isInteger(personId) || personId < 0) {
		return json({ msg: 'invalid person id' }, 400);
	}

	if (personId !== userId) {
		const admin = await client.GET('/admin/{id1}/{id2}', {
			params: {
				path: { id1: personId, id2: userId },
				header: { 'X-User-ID': userId }
			}
		});
		if (!admin.response.ok) {
			return json({ msg: 'user can not manage this person' }, 403);
		}
	}

	if (Number(event.request.headers.get('content-length')) > MAX_MEDIA_SIZE + 1024 * 1024) {
		return json({ msg: 'file too large' }, 413);
	}

	let file: FormDataEntryValue | null;
	try {
		file = (await event.request.formData()).get('file');
	} catch {
		return json({ msg: 'invalid multipart body' }, 400);
	}

	if (!(file instanceof File) || file.size === 0) {
		return json({ msg: 'file is missing' }, 400);
	}
	if (!ALLOWED_TYPE.test(file.type)) {
		return json({ msg: 'unsupported media type' }, 415);
	}
	if (file.size > MAX_MEDIA_SIZE) {
		return json({ msg: 'file too large' }, 413);
	}

	const extension = /\.([a-z0-9]{1,8})$/i.exec(file.name)?.[1].toLowerCase() ?? 'bin';
	const key = `people/${personId}/${uuidv4()}.${extension}`;
	try {
		await bucket.put(key, file, { httpMetadata: { contentType: file.type } });
	} catch (error) {
		console.error('Error storing media', error);
		return json({ msg: 'failed to store media' }, 500);
	}

	return json({ url: `/api/media/${key}`, key }, 201);
}
