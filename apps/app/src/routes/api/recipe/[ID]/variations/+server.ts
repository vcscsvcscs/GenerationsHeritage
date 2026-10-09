import { client } from '$lib/api/client';
import { toResponse, userHeader, withSessionAndId } from '$lib/server/proxy';
import type { RequestEvent } from './$types';

export const GET = withSessionAndId(async (_event: RequestEvent, userId, id) =>
	toResponse(
		await client.GET('/recipe/{id}/variations', {
			params: { path: { id }, header: userHeader(userId) }
		})
	)
);
