import { redirect } from '@sveltejs/kit';
import { parseFamilyTree } from '$lib/graph/fetch_family_tree';
import type { components } from '$lib/api/api.gen';
import type { RequestEvent } from './$types';
import { browser } from '$app/environment';
import type { Layout } from '$lib/graph/model';

export async function load(event: RequestEvent) {
	if (event.locals.session === null /*|| event.locals.familytree === nul*/) {
		return redirect(302, '/login');
	}

	//prevent loading in developer mode, due to some issues with universal load, even if this is a server only ts,it will still run on client in dev mode idk
	if (browser) {
		return {};
	}

	const response = await event.fetch('/api/family_tree?with_out_spouse=false', {
		method: 'GET'
	});

	if (response.status !== 200) {
		throw new Error(await response.text());
	}

	const data = (await response.json()) as components['schemas']['FamilyTree'];

	const layout = parseFamilyTree(data) as Layout & { id: string };
	layout.id = event.locals.session.userId;

	return layout;
}
