import { google } from '$lib/server/oauth';
import { ObjectParser } from '@pilcrowjs/object-parser';
import { browser } from '$app/environment';
import { client } from '$lib/api/client';
import { type components } from '$lib/api/api.gen';
import { createSession, generateSessionToken, setSessionTokenCookie } from '$lib/server/session';
import { decodeIdToken } from 'arctic';
import {
	missing_field,
	last_name,
	first_name,
	mothers_first_name,
	mothers_last_name,
	born,
	failed_to_create_user,
	biological_sex
} from '$lib/paraglide/messages';

import type { PageServerLoad, Actions, RequestEvent } from './$types';
import type { OAuth2Tokens } from 'arctic';
import type { PersonProperties } from '$lib/model';
import { error, redirect, fail } from '@sveltejs/kit';

const StorageLimit = 200 * 1024 * 1024;

export const load: PageServerLoad = async (event: RequestEvent) => {
	//prevent loading in developer mode, due to some issues with universal load, even if this is a server only ts,it will still run on client in dev mode idk
	if (browser) {
		return {};
	}

	let already_loaded = event.cookies.get('already_loaded') ?? null;
	if (already_loaded !== null) {
		return {};
	}


	const storedState = event.cookies.get('google_oauth_state') ?? null;
	const codeVerifier = event.cookies.get('google_code_verifier') ?? null;
	const code = event.url.searchParams.get('code');
	const state = event.url.searchParams.get('state');

	if (storedState === null || codeVerifier === null || code === null || state === null) {
		return error(400, { message: 'Please restart the process.' });
	}
	if (storedState !== state) {
		return error(400, { message: 'Please restart the process.' });
	}

	let tokens: OAuth2Tokens;
	try {
		tokens = await google.validateAuthorizationCode(code, codeVerifier);
	} catch (e) {
		return error(400, { message: 'Failed to validate authorization code with ' + e });
	}

	const claims = decodeIdToken(tokens.idToken());
	const claimsParser = new ObjectParser(claims);

	const sub = claimsParser.getString('sub');
	const family_name = claimsParser.getString('family_name');
	const first_name = claimsParser.getString('given_name');
	const email = claimsParser.getString('email');

	const response = await client.GET('/person/google/{google_id}', {
		params: {
			path: { google_id: sub }
		}
	});

	if (response.response.status === 200) {
		if (response.data?.Id) {
			if (!event.platform || !event.platform.env || !event.platform.env.GH_SESSIONS) {
				return error(500, {
					message: 'Server configuration error. GH_SESSIONS KeyValue store missing'
				});
			}

			const sessionToken = generateSessionToken(String(response.data.Id));
			const session = await createSession(
				sessionToken,
				response.data.Id,
				event.platform.env.GH_SESSIONS
			);
			if (session === null) {
				return error(500, {
					message: 'Failed to create session'
				});
			}

			setSessionTokenCookie(event, sessionToken, session.expiresAt);

			return redirect(302, '/');
		}
	}

	let personP: PersonProperties = {
		google_id: sub,
		first_name: first_name,
		last_name: family_name,
		email: email
	};

	event.cookies.set('already_loaded', 'true',{
		path: '/login/google/callback',
		sameSite: 'lax',
		httpOnly: true,
		maxAge: 60 * 10,
		secure: import.meta.env.PROD,
	})

	return {
		props: personP
	};
};

export const actions: Actions = {
	register: register
};

async function register(event: RequestEvent) {
	if (browser) {
		return {};
	}

	const data = await event.request.formData();
	let parsedData: components['schemas']['PersonRegistration'] = {
		first_name: data.get('first_name'),
		last_name: data.get('last_name'),
		email: data.get('email'),
		biological_sex: data.get('biological_sex'),
		born: data.get('birth_date'),
		mothers_first_name: data.get('mothers_first_name'),
		mothers_last_name: data.get('mothers_last_name'),
		google_id: data.get('google_id'),
		limit: StorageLimit,
	} as components['schemas']['PersonRegistration'];

	if (!event.platform || !event.platform.env || !event.platform.env.GH_SESSIONS) {
		return fail(500, { data: parsedData, message: 'Server configuration error. GH_SESSIONS KeyValue store missing' });
	}

	const first_name_f = data.get('first_name');
	if (first_name_f === null || first_name_f === '') {
		return fail(400, {
			data:parsedData,
			message: missing_field({
				field: first_name()
			})
		});
	}

	const google_id = data.get('google_id');
	if (google_id === null || google_id === '') {
		return fail(400, {
			data: parsedData,
			message: missing_field({
				field: 'google_id'
			})
		});
	}

	const last_name_f = data.get('last_name');
	if (last_name_f === null || last_name_f === '') {
		return fail(400, {
			data: parsedData,
			message: missing_field({
				field: last_name()
			})
		});
	}

	const email = data.get('email');
	if (email === null || email === '') {
		return fail(400, {
			message: missing_field({
				field: 'Email'
			})
		});
	}
	const birth_date = data.get('birth_date');
	if (birth_date === null || birth_date === '') {
		return fail(400, {
			message: missing_field({
				field: born()
			})
		});
	}

	const bbiological_sex = data.get('biological_sex');
	if (bbiological_sex === null || bbiological_sex === '') {
		return fail(400, {
			data: parsedData,
			message: missing_field({
				field: biological_sex()
			})
		});
	} else if (
		!['male', 'female', 'intersex', 'unknown', 'other'].includes(bbiological_sex.toString())
	) {
		return fail(400, {
			data: parsedData,
			message: `Invalid value for biological_sex. Must be one of "male", "female", "intersex", "unknown", or "other".`
		});
	}

	const mothers_first_name_f = data.get('mothers_first_name');
	if (mothers_first_name_f === null || mothers_first_name_f === '') {
		return fail(400, {
			data: parsedData,
			message: missing_field({
				field: mothers_first_name()
			})
		});
	}
	const mothers_last_name_f = data.get('mothers_last_name');
	if (mothers_last_name_f === null) {
		return fail(400, {
			data: parsedData,
			message: missing_field({
				field: mothers_last_name()
			})
		});
	}

	const parsed_date = new Date(birth_date as string);
	let personP: components['schemas']['PersonRegistration'] = {
		first_name: first_name_f as string,
		last_name: last_name_f as string,
		email: email as string,
		born: parsed_date.toISOString().split('T')[0],
		mothers_first_name: mothers_first_name_f as string,
		mothers_last_name: mothers_last_name_f as string,
		biological_sex:
			bbiological_sex as components['schemas']['PersonRegistration']['biological_sex'],
		limit: StorageLimit
	};

	let response = await client.POST('/person/google/{google_id}', {
		params: {
			data: parsedData,
			path: { google_id: google_id.toString() }
		},
		body: personP
	});

	if (response.response.status !== 200) {
		return fail(400, {
			data: parsedData,
			message: failed_to_create_user() + response.error?.msg
		});
	}

	if (response.data === undefined) {
		return fail(400, {
			data: parsedData,
			message: failed_to_create_user() + 'No user data returned'
		});
	}

	if (response.data.Id === undefined) {
		return fail(400, {
			data: parsedData,
			message: failed_to_create_user() + 'No user ID returned'
		});
	}

	if (!event.platform) {
		return fail(500, {
			data: parsedData,
			message: 'Server configuration error. GH_SESSIONS KeyValue store missing'
		});
	}

	const sessionToken = generateSessionToken(String(response.data.Id));
	const session = await createSession(
		sessionToken,
		response.data.Id,
		event.platform.env.GH_SESSIONS
	);
	if (session === null) {
		return fail(500, {
			data: parsedData,
			message: failed_to_create_user() + 'Failed to create session'
		});
	}

	setSessionTokenCookie(event, sessionToken, session.expiresAt);

	event.cookies.delete('already_loaded',
		{
			path: '/login/google/callback',
			sameSite: 'lax',
			httpOnly: true,
			maxAge: 0,
			secure: import.meta.env.PROD
		}
	);
	return redirect(302, '/');
}
