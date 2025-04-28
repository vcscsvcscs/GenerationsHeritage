import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export async function POST(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	let message = (await event.request.json()) as components['schemas']['Message'];
	message.edited = null;
	message.sent_at = new Date(Date.now()).toISOString();

	const response = await client.POST('/comment/{id}', {
		params: {
			path: { id: Number(event.params.ID) },
			header: { 'X-User-ID': event.locals.session.userId }
		},
		body: message
	});

	return new Response(await response.response.json(), {
		status: response.response.status
	});
}

export async function GET(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.GET('/comment/{id}', {
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

export async function DELETE(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	const response = await client.DELETE('/comment/{id}', {
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

export async function PATCH(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	let message = (await event.request.json()) as components['schemas']['Message'];
	message.edited = new Date(Date.now()).toISOString();

	const response = await client.PATCH('/comment/{id}', {
		params: {
			path: { id: Number(event.params.ID) },
			header: { 'X-User-ID': event.locals.session.userId }
		},
		body: message
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
