import { client } from '$lib/api/client';
import { badRequest, readJson, toResponse, userHeader, withSession } from '$lib/server/proxy';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

const personIdOf = (event: RequestEvent, userId: number) =>
	event.params.ID === 'me' ? userId : Number(event.params.ID);

export const GET = withSession(async (event: RequestEvent, userId) => {
	const personId = personIdOf(event, userId);
	if (!Number.isInteger(personId)) return badRequest('Invalid person ID');

	return toResponse(
		await client.GET('/person/{id}/recipes', {
			params: { path: { id: personId }, header: userHeader(userId) }
		})
	);
});

export const POST = withSession(async (event: RequestEvent, userId) => {
	const personId = personIdOf(event, userId);
	if (!Number.isInteger(personId)) return badRequest('Invalid person ID');

	const body = await readJson<{
		recipe: components['schemas']['RecipeProperties'];
		relationship?: components['schemas']['LikesProperties'];
	}>(event.request);
	if (body?.recipe == null) return badRequest('Missing recipe');

	return toResponse(
		await client.POST('/person/{id}/recipes', {
			params: { path: { id: personId }, header: userHeader(userId) },
			body
		})
	);
});
