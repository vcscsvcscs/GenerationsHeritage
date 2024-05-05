import { PUBLIC_API_URL } from '$env/static/public';
import { user } from '$lib/stores';

let auth_token: string;

user.subscribe((value) => {
	if (value) {
		auth_token = value.access_token;
	}
});

async function fetch_family_tree() {
	console.log(PUBLIC_API_URL);

	const response = await fetch(PUBLIC_API_URL + '/familyTree?id=8a8b9b05bdc24550a5cc73e0b55e8d7d', {
		method: 'GET',
		headers: {
			Accept: 'application/json',
			Authorization: auth_token
		}
	});
	const data = await response.json();
	return data;
}

export { fetch_family_tree };
