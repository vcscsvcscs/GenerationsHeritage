import { error, redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';

export async function DELETE(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.DELETE('/person/{id}/hard-delete', {
		params: {
			path: { id: Number(event.params.ID) },
			header: { 'X-User-ID': event.locals.session.userId }
		}
	});

	if (response.response.ok) {
		return new Response(null, {
			status: response.response.status
		});
	} else {
		return new Response(JSON.stringify(response.error), {
			status: response.response.status
		});
	}
}
