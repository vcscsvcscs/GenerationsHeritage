import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export async function GET(event: RequestEvent): Promise<Response> {
    if (event.locals.session === null) {
        return redirect(302, '/login');
    }

    const response = await client.GET(
        '/person/{id}',
        {
            params: {
                path: { id: Number(event.params.ID) },
                header: { 'X-User-ID': event.locals.session.userId }
            }
        }
    );

    return new Response(await response.response.json(), {
        status: response.response.status,
    });
}

export async function DELETE(event: RequestEvent): Promise<Response> {
    if (event.locals.session === null) {
        return redirect(302, '/login');
    }

    const response = await client.DELETE(
        '/person/{id}',
        {
            params: {
                path: { id: Number(event.params.ID) },
                header: { 'X-User-ID': event.locals.session.userId }
            }
        }
    );

    return new Response(await response.response.json(), {
        status: response.response.status,
    });
}

export async function PATCH(event: RequestEvent): Promise<Response> {
    if (event.locals.session === null) {
        return redirect(302, '/login');
    }

    const response = await client.PATCH(
        '/person/{id}',
        {
            params: {
                path: { id: Number(event.params.ID) },
                header: { 'X-User-ID': event.locals.session.userId }
            },
            body: await event.request.json() as components['schemas']['PersonProperties']
        }
    );

    return new Response(await response.response.json(), {
        status: response.response.status,
    });
}