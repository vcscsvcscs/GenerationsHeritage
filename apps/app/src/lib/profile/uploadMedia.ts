import { file_too_large, unsupported_file_type, upload_failed } from '$lib/paraglide/messages';

export async function uploadMedia(personId: string | number, file: File): Promise<string> {
	const body = new FormData();
	body.append('file', file);

	let response: Response;
	try {
		response = await fetch(`/api/person/${personId}/media`, { method: 'POST', body });
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
