import { redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';

export async function GET(event: RequestEvent): Promise<Response> {
    if (event.locals.session === null) {
        return redirect(302, '/login');
    }

    const response = await client.GET(
        '/relationship/{id1}/{id2}',
        {
            params: {
                path: { id1: Number(event.params.ID1), id2: Number(event.params.ID2) },
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
        '/relationship/{id1}/{id2}',
        {
            params: {
                path: { id1: Number(event.params.ID1), id2: Number(event.params.ID2) },
                header: { 'X-User-ID': event.locals.session.userId }
            },
            body: {
                relationship: event.request.json() as components['schemas']['FamilyRelationship']
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
        '/relationship/{id1}/{id2}',
        {
            params: {
                path: { id1: Number(event.params.ID1), id2: Number(event.params.ID2) },
                header: { 'X-User-ID': event.locals.session.userId }
            }
        }
    );

    return new Response(await response.response.json(), {
        status: response.response.status,
    });
}