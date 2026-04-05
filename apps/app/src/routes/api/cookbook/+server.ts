import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';

export async function GET(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const distance = Number(event.url.searchParams.get('distance') ?? 3);

	const response = await client.GET('/cookbook', {
		params: {
			query: { distance },
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
