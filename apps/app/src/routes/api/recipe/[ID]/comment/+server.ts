import { client } from '$lib/api/client';
import { badRequest, readJson, toResponse, userHeader, withSessionAndId } from '$lib/server/proxy';
import type { RequestEvent } from './$types';

async function readMessage(request: Request): Promise<string | null> {
	const body = await readJson<{ message?: unknown }>(request);
	return typeof body?.message === 'string' && body.message.trim() !== '' ? body.message : null;
}

export const GET = withSessionAndId(async (_event: RequestEvent, userId, id) =>
	toResponse(
		await client.GET('/recipe/{id}/comment', {
			params: { path: { id }, header: userHeader(userId) }
		})
	)
);

export const POST = withSessionAndId(async (event: RequestEvent, userId, id) => {
	const message = await readMessage(event.request);
	if (message === null) return badRequest('Message must not be empty');

	return toResponse(
		await client.POST('/recipe/{id}/comment', {
			params: { path: { id }, header: userHeader(userId) },
			body: { message }
		})
	);
});

export const PATCH = withSessionAndId(async (event: RequestEvent, userId, id) => {
	const message = await readMessage(event.request);
	if (message === null) return badRequest('Message must not be empty');

	return toResponse(
		await client.PATCH('/recipe/{id}/comment', {
			params: { path: { id }, header: userHeader(userId) },
			body: { message }
		})
	);
});

export const DELETE = withSessionAndId(async (_event: RequestEvent, userId, id) =>
	toResponse(
		await client.DELETE('/recipe/{id}/comment', {
			params: { path: { id }, header: userHeader(userId) }
		})
	)
);
