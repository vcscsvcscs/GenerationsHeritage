import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export async function PATCH(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.PATCH('/recipe/{id}', {
		params: {
			path: { id: Number(event.params.ID) },
			header: { 'X-User-ID': event.locals.session.userId }
		},
		body: (await event.request.json()) as components['schemas']['RecipeProperties']
	});

	if (response.response.ok) {
		return new Response(JSON.stringify(response.data), {
			status: response.response.status
		});
	} else {
		return new Response(JSON.stringify(response.error), {
			status: response.response.status
		});
	}
}

export async function DELETE(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.DELETE('/recipe/{id}', {
		params: {
			path: { id: Number(event.params.ID) },
			header: { 'X-User-ID': event.locals.session.userId }
		}
	});

	if (response.response.ok) {
		return new Response(JSON.stringify(response.data), {
			status: response.response.status
		});
	} else {
		return new Response(JSON.stringify(response.error), {
			status: response.response.status
		});
	}
}
