import { client } from '$lib/api/client';
import { badRequest, readJson, toResponse, userHeader, withSessionAndId } from '$lib/server/proxy';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export const POST = withSessionAndId(async (event: RequestEvent, userId, id) => {
	const body = await readJson<{
		recipe: components['schemas']['RecipeProperties'];
		variation_notes?: string | null;
		relationship?: components['schemas']['LikesProperties'];
	}>(event.request);
	if (body?.recipe == null) return badRequest('Missing recipe');

	return toResponse(
		await client.POST('/recipe/{id}/variation', {
			params: { path: { id }, header: userHeader(userId) },
			body
		})
	);
});
