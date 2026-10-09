import { client } from '$lib/api/client';
import { badRequest, toResponse, userHeader, withSession } from '$lib/server/proxy';
import type { RequestEvent } from './$types';

export const GET = withSession(async (event: RequestEvent, userId) => {
	const distance = Number(event.url.searchParams.get('distance') ?? 3);
	if (!Number.isInteger(distance) || distance < 0) return badRequest('Invalid distance');

	return toResponse(
		await client.GET('/cookbook', {
			params: { query: { distance }, header: userHeader(userId) }
		})
	);
});
