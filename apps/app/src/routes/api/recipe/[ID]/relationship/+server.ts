import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';

export async function POST(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const body = await event.request.json();
	// Auto-fill person id from session if not provided or 0
	if (!body.id) {
		body.id = event.locals.session.userId;
	}

	const response = await client.POST('/recipe/{id}/relationship', {
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

export async function DELETE(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const personIdParam = event.url.searchParams.get('personId');
	const personId = personIdParam === 'me' ? event.locals.session.userId : Number(personIdParam);

	const response = await client.DELETE('/recipe/{id}/relationship', {
		params: {
			path: { id: Number(event.params.ID) },
			query: { personId },
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
