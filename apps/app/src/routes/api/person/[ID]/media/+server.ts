import { redirect } from '@sveltejs/kit';
import { v4 as uuidv4 } from 'uuid';
import { MEDIA_KEY_PATTERN, json, requireManage } from '$lib/server/media';
import type { RequestEvent } from './$types';

const MAX_MEDIA_SIZE = 25 * 1024 * 1024;

const ALLOWED_TYPE = /^(image|video|audio)\/(?!svg)[\w.+-]+$/;

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

	const kind = event.url.searchParams.get('kind');
	if (kind !== null && kind !== 'profile_picture') {
		return json({ msg: 'invalid kind' }, 400);
	}

	const denied = await requireManage(personId, userId);
	if (denied) {
		return denied;
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
	if (
		!ALLOWED_TYPE.test(file.type) ||
		(kind === 'profile_picture' && !/^image\//.test(file.type))
	) {
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

export async function DELETE(event: RequestEvent): Promise<Response> {
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
	const key = event.url.searchParams.get('key') ?? '';
	if (!Number.isInteger(personId) || personId < 0) {
		return json({ msg: 'invalid person id' }, 400);
	}
	if (MEDIA_KEY_PATTERN.exec(key)?.[1] !== String(personId)) {
		return json({ msg: 'invalid media key' }, 400);
	}

	const denied = await requireManage(personId, event.locals.session.userId);
	if (denied) {
		return denied;
	}

	try {
		await bucket.delete(key);
	} catch (error) {
		console.error('Error deleting media', error);
		return json({ msg: 'failed to delete media' }, 500);
	}

	return new Response(null, { status: 204 });
}
