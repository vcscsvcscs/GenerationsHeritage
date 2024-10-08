import { PUBLIC_API_URL } from '$env/static/public';
import { user } from '$lib/stores';

let auth_token: string;

user.subscribe((value) => {
	if (value) {
		auth_token = value.access_token;
	}
});

async function fetch_family_tree() {
	const response = await fetch(PUBLIC_API_URL + '/familyTree', {
		method: 'GET',
		headers: {
			Accept: 'application/json',
			Authorization: 'Bearer ' + auth_token
		}
	});
	const data = await response.json();
	return data;
}

export { fetch_family_tree };
