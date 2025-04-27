import { error, redirect } from '@sveltejs/kit';
import { client } from '$lib/api/client';
import type { RequestEvent } from './$types';
import type { components } from '$lib/api/api.gen';


async function GET(event: RequestEvent): Promise<Response> {
	if (event.locals.session === null /*|| event.locals.familytree === nul*/) {
		return redirect(302, '/login');
	}

	const response = await client.GET(event.url.searchParams.get('with_out_spouse') === 'true' ? '/family-tree-with-spouses':'/family-tree', {
		params: {
			header: { 'X-User-ID': event.locals.session.userId }
		}
	});

	if (response.response.status !== 200) {
		return error(500, {
			message: response.error?.msg || 'Failed to fetch family tree'
		});
	}

	if (response.data === null || response.data?.people === null || response.data?.people === undefined || response.data?.people.length === 0) {
		return error(500, {
			message: 'Family tree is empty'
		});
	}

	var peopleToReturn : Array<components['schemas']['OptimizedPersonNode']> = [];
	for (const person of response.data.people) {
		let newPerson = person

		if (newPerson.profile_picture!== null && newPerson.profile_picture !== undefined) {
		
		}

		peopleToReturn.push(newPerson);
	}

	return new Response(JSON.stringify(peopleToReturn), {
		status: 200,
	});
}
