import { client } from '$lib/api/client';
import { badRequest, readJson, toResponse, userHeader, withSessionAndId } from '$lib/server/proxy';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export const GET = withSessionAndId(async (_event: RequestEvent, userId, id) =>
	toResponse(
		await client.GET('/recipe/{id}', {
			params: { path: { id }, header: userHeader(userId) }
		})
	)
);

export const PATCH = withSessionAndId(async (event: RequestEvent, userId, id) => {
	const body = await readJson<components['schemas']['RecipeProperties']>(event.request);
	if (body === null) return badRequest('Invalid JSON body');

	return toResponse(
		await client.PATCH('/recipe/{id}', {
			params: { path: { id }, header: userHeader(userId) },
			body
		})
	);
});

export const DELETE = withSessionAndId(async (_event: RequestEvent, userId, id) =>
	toResponse(
		await client.DELETE('/recipe/{id}', {
			params: { path: { id }, header: userHeader(userId) }
		})
	)
);
