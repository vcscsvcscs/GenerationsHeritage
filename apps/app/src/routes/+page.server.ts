import { fail, redirect } from '@sveltejs/kit';
import { deleteSessionTokenCookie, invalidateSession } from '$lib/server/session';
import { client } from '$lib/api/client';
import type { Actions, RequestEvent } from './$types';
import { browser } from '$app/environment';

export async function load(event: RequestEvent) {
	if (event.locals.session === null /*|| event.locals.familytree === nul*/) {
		return redirect(302, '/login');
	}

	//prevent loading in developer mode, due to some issues with universal load, even if this is a server only ts,it will still run on client in dev mode idk
	if (browser) {
		return {};
	}

	client.GET('/family-tree-with-spouses', {
		params: {
			header: { "X-User-ID": event.locals.session },
		}
	}).then((response) => {

		if (response.response.status === 200) {
			return response.data;
		} else {
			return fail(response.response.status, { message: response.error?.msg || 'An error occurred' });
		}
	});
}

export const actions: Actions = {
	logout: logout
};

async function logout(event: RequestEvent) {
	if (event.locals.session === null) {
		return fail(401);
	}

	if (event.platform && event.platform.env && event.platform.env.GH_SESSIONS) {
		invalidateSession(event.locals.session.id, event.platform.env.GH_SESSIONS);
	} else {
		return fail(500, { message: 'Server configuration error' });
	}

	deleteSessionTokenCookie(event);

	return redirect(302, '/login');
}
