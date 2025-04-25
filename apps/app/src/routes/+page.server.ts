import { fail, redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type {  RequestEvent } from './$types';
import { browser } from '$app/environment';

export async function load(event: RequestEvent) {
	if (event.locals.session === null /*|| event.locals.familytree === nul*/) {
		return redirect(302, '/login');
	}

	//prevent loading in developer mode, due to some issues with universal load, even if this is a server only ts,it will still run on client in dev mode idk
	if (browser) {
		return {};
	}

	const response = await client
		.GET('/family-tree-with-spouses', {
			params: {
				header: { 'X-User-ID': event.locals.session.userId },
			}
		})

	if (response.response.status === 200) {
		return response.data;
	} else {
		return fail(response.response.status, {
			message: response.error?.msg || 'An error occurred'
		});
	}

}
