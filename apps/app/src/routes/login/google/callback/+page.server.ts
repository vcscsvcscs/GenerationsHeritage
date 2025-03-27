import { google } from '$lib/server/oauth';
import { ObjectParser } from '@pilcrowjs/object-parser';
import { createUser, getUserFromGoogleId } from '$lib/server/user';
import { DB } from '$lib/server/db';
import { browser } from '$app/environment';
import { Date as neoDate } from 'neo4j-driver';
import { createSession, generateSessionToken, setSessionTokenCookie } from '$lib/server/session';
import { decodeIdToken } from 'arctic';
import {
	missing_field,
	last_name,
	first_name,
	mothers_first_name,
	mothers_last_name,
	born,
	failed_to_create_user
} from '$lib/paraglide/messages';

import type { PageServerLoad, Actions, RequestEvent, PageData } from './$types';
import type { OAuth2Tokens } from 'arctic';
import type { PersonProperties } from '$lib/model';
import { error, redirect, fail } from '@sveltejs/kit';

const StorageLimit = 200 * 1024 * 1024;

export const load: PageServerLoad = async (event: RequestEvent) => {
	//prevent loading in developer mode, due to some issues with universal load, even if this is a server only ts,it will still run on client in dev mode idk
	if (browser) {
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

	//	if (!event.platform || !event.platform.env || !event.platform.env.GH_SESSIONS) {
	//		return error(500, { message: "Server configuration error. GH_SESSIONS KeyValue store missing" });
	//	}

	const claims = decodeIdToken(tokens.idToken());
	const claimsParser = new ObjectParser(claims);

	const sub = claimsParser.getString('sub');
	const family_name = claimsParser.getString('family_name');
	const first_name = claimsParser.getString('given_name');
	const email = claimsParser.getString('email');

	const dbSession = DB.session();
	const existingUser = await getUserFromGoogleId(dbSession, sub);
	dbSession.close();

	let eUser = existingUser.records.pop();
	if (eUser !== null && eUser?.get('elementId') !== undefined) {
		const sessionToken = generateSessionToken();
		//		const session = await createSession(sessionToken, eUser.get('elementId'), event.platform.env.GH_SESSIONS);
		//		setSessionTokenCookie(event, sessionToken, session.expiresAt);

		return redirect(302, '/');
	}

	let personP: PersonProperties = {
		google_id: sub,
		first_name: first_name,
		last_name: family_name,
		email: email,
		allow_admin_access: false,
		limit: StorageLimit,
		verified: false
	};

	return {
		props: personP
	};
};

export const actions: Actions = {
	register: register
};

async function register(event: RequestEvent) {
	const data = await event.request.formData();
	//	if (!event.platform || !event.platform.env || !event.platform.env.GH_SESSIONS) {
	//		return fail(500, { message: "Server configuration error. GH_SESSIONS KeyValue store missing" });
	//	}
	const google_id = data.get('google_id');
	if (google_id === null) {
		return fail(400, {
			message: missing_field({
				field: 'google_id'
			})
		});
	}
	const first_name_f = data.get('first_name');
	if (first_name_f === null) {
		return fail(400, {
			message: missing_field({
				field: first_name()
			})
		});
	}
	const last_name_f = data.get('last_name');
	if (last_name_f === null) {
		return fail(400, {
			message: missing_field({
				field: last_name()
			})
		});
	}
	const email = data.get('email');
	if (email === null) {
		return fail(400, {
			message: missing_field({
				field: 'Email'
			})
		});
	}
	const birth_date = data.get('birth_date');
	if (birth_date === null) {
		return fail(400, {
			message: missing_field({
				field: born()
			})
		});
	}
	const mothers_first_name_f = data.get('mothers_first_name');
	if (mothers_first_name_f === null) {
		return fail(400, {
			message: missing_field({
				field: mothers_first_name()
			})
		});
	}
	const mothers_last_name_f = data.get('mothers_last_name');
	if (mothers_last_name_f === null) {
		return fail(400, {
			message: missing_field({
				field: mothers_last_name()
			})
		});
	}

	const parsed_date = new Date(birth_date as string);

	let personP: PersonProperties = {
		google_id: google_id as string,
		first_name: first_name_f as string,
		last_name: last_name_f as string,
		email: email as string,
		born: new neoDate(
			parsed_date.getFullYear(),
			parsed_date.getUTCMonth(),
			parsed_date.getUTCDate()
		),
		mothers_first_name: mothers_first_name_f as string,
		mothers_last_name: mothers_last_name_f as string,
		allow_admin_access: false,
		limit: StorageLimit,
		verified: false
	};

	const dbSession = DB.session();
	const user = (await createUser(dbSession, personP)).records.pop();
	if (user === null || user === undefined) {
		dbSession.close();

		return fail(500, { message: failed_to_create_user() });
	}
	dbSession.close();

	const sessionToken = generateSessionToken();
	//	const session = await createSession(sessionToken, user.get('elementId'), event.platform.env.GH_SESSIONS);
	//	setSessionTokenCookie(event, sessionToken, session.expiresAt);

	return redirect(302, '/');
}
