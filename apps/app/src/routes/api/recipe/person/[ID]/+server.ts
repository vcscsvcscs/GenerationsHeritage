import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export async function GET(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.GET('/person/{id}/recipes', {
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

export async function POST(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const body = (await event.request.json()) as {
		recipe: components['schemas']['RecipeProperties'];
		relationship?: components['schemas']['LikesProperties'];
	};

	const response = await client.POST('/person/{id}/recipes', {
		params: {
			path: { id: Number(event.params.ID) },
			header: { 'X-User-ID': event.locals.session.userId }
		},
		body
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
