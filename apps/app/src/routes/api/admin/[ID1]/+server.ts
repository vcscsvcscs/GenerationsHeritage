import { client } from '$lib/api/client';
import { redirect } from '@sveltejs/kit';
import type { RequestEvent } from './$types';

export async function GET(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.GET('/admin/{id1}', {
		params: {
			path: { id1: Number(event.params.ID1) },
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
