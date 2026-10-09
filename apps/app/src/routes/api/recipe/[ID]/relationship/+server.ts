import { client } from '$lib/api/client';
import { badRequest, readJson, toResponse, userHeader, withSessionAndId } from '$lib/server/proxy';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export const POST = withSessionAndId(async (event: RequestEvent, userId, id) => {
	const body = await readJson<components['schemas']['RecipeRelationshipInput']>(event.request);
	if (body === null) return badRequest('Invalid JSON body');

	return toResponse(
		await client.POST('/recipe/{id}/relationship', {
			params: { path: { id }, header: userHeader(userId) },
			body
		})
	);
});

export const DELETE = withSessionAndId(async (event: RequestEvent, userId, id) => {
	const personIdParam = event.url.searchParams.get('personId');
	const personId = personIdParam === 'me' ? userId : Number(personIdParam);
	if (personIdParam === null || !Number.isInteger(personId)) return badRequest('Invalid personId');

	return toResponse(
		await client.DELETE('/recipe/{id}/relationship', {
			params: { path: { id }, query: { personId }, header: userHeader(userId) }
		})
	);
});
