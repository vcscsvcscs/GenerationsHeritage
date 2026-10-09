import { client } from '$lib/api/client';

export const MEDIA_KEY_PATTERN = /^people\/(\d+)\/[0-9a-f-]{36}\.[a-z0-9]{1,8}$/;

export function json(body: unknown, status: number): Response {
	return new Response(JSON.stringify(body), { status });
}

/**
 * Mirrors the db-adapter's CouldManagePersonUnknownAdmin: the owner always passes,
 * anyone else needs an Admin relationship. Returns an error response, or null when allowed.
 */
export async function requireManage(personId: number, userId: number): Promise<Response | null> {
	if (personId === userId) {
		return null;
	}

	let status: number;
	try {
		const admin = await client.GET('/admin/{id1}/{id2}', {
			params: {
				path: { id1: personId, id2: userId },
				header: { 'X-User-ID': userId }
			}
		});
		status = admin.response.status;
	} catch (error) {
		console.error('Error checking admin relationship', error);
		return json({ msg: 'permission check unavailable' }, 502);
	}

	if (status >= 200 && status < 300) {
		return null;
	}
	if (status === 401 || status === 403) {
		return json({ msg: 'user can not manage this person' }, 403);
	}
	if (status === 404) {
		return json({ msg: 'person not found' }, 404);
	}
	return json({ msg: 'permission check failed' }, 502);
}
