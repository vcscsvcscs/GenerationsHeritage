import { error, redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export async function POST(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.POST('/person', {
		params: {
			header: { 'X-User-ID': event.locals.session.userId }
		},
		body: (await event.request.json()) as components['schemas']['PersonRegistration']
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
