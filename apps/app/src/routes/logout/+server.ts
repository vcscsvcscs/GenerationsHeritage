import { error, redirect } from '@sveltejs/kit';
import { invalidateSession, deleteSessionTokenCookie } from '$lib/server/session';

import type { RequestEvent } from './$types';

export async function GET(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null) {
		return redirect(302, '/login');
	}

	if (event.platform && event.platform.env && event.platform.env.GH_SESSIONS) {
		await invalidateSession(event.locals.session.id, event.platform.env.GH_SESSIONS);
	} else {
		return error(500, { message: 'Server configuration error' });
	}

	deleteSessionTokenCookie(event);

	return redirect(302, '/login');
}
