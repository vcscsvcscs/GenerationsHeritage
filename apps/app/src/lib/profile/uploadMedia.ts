import type { components } from '$lib/api/api.gen';
import { file_too_large, unsupported_file_type, upload_failed } from '$lib/paraglide/messages';

const MEDIA_URL_PREFIX = '/api/media/';

export async function uploadMedia(
	personId: string | number,
	file: File,
	kind?: 'profile_picture'
): Promise<string> {
	const body = new FormData();
	body.append('file', file);

	let response: Response;
	try {
		response = await fetch(`/api/person/${personId}/media${kind ? `?kind=${kind}` : ''}`, {
			method: 'POST',
			body
		});
	} catch (e) {
		console.error('Error uploading media', e);
		throw new Error(upload_failed());
	}

	if (!response.ok) {
		throw new Error(
			response.status === 413
				? file_too_large()
				: response.status === 415
					? unsupported_file_type()
					: upload_failed()
		);
	}

	return ((await response.json()) as { url: string }).url;
}

export async function deleteMedia(personId: string | number, url: string): Promise<void> {
	if (!url.startsWith(MEDIA_URL_PREFIX)) return;

	const key = encodeURIComponent(url.slice(MEDIA_URL_PREFIX.length));
	try {
		const response = await fetch(`/api/person/${personId}/media?key=${key}`, {
			method: 'DELETE'
		});
		if (!response.ok) console.error('Error deleting media, status:', response.status);
	} catch (e) {
		console.error('Error deleting media', e);
	}
}

export async function deleteUnreferencedMedia(
	personId: string | number,
	urls: Iterable<string>,
	person: components['schemas']['PersonProperties']
): Promise<void> {
	const referenced = new Set([
		person.profile_picture,
		...[person.photos, person.videos, person.audios].flatMap((items) =>
			(items ?? []).map((item) => item.url)
		)
	]);
	await Promise.all(
		[...new Set(urls)]
			.filter((url) => !referenced.has(url))
			.map((url) => deleteMedia(personId, url))
	);
}
