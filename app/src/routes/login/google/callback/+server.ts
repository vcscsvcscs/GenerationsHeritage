import { fail } from "@sveltejs/kit";
import { google } from "$lib/server/oauth";
import { ObjectParser } from "@pilcrowjs/object-parser";
import { createUser, getUserFromGoogleId } from "$lib/server/user";
import { driverInstance } from "$lib/server/db";
import { createSession, generateSessionToken, setSessionTokenCookie } from "$lib/server/session";
import { decodeIdToken } from "arctic";

import type { RequestEvent } from "./$types";
import type { OAuth2Tokens } from "arctic";
import { middle_name } from "$lib/paraglide/messages";

export async function GET(event: RequestEvent): Promise<Response> {
	const storedState = event.cookies.get("google_oauth_state") ?? null;
	const codeVerifier = event.cookies.get("google_code_verifier") ?? null;
	const code = event.url.searchParams.get("code");
	const state = event.url.searchParams.get("state");

	if (storedState === null || codeVerifier === null || code === null || state === null) {
		return new Response("Please restart the process.", {
			status: 400
		});
	}
	if (storedState !== state) {
		return new Response("Please restart the process.", {
			status: 400
		});
	}

	let tokens: OAuth2Tokens;
	try {
		tokens = await google.validateAuthorizationCode(code, codeVerifier);
	} catch (e) {
		return new Response("Please restart the process.", {
			status: 400
		});
	}

	if (!event.platform || !event.platform.env || !event.platform.env.GH_SESSIONS) {
		return new Response("Server configuration error. GH_SESSIONS KeyValue store missing", {
			status: 500
		});
	}

	const claims = decodeIdToken(tokens.idToken());
	const claimsParser = new ObjectParser(claims);

	const googleId = claimsParser.getString("sub");
	const family_name = claimsParser.getString("family_name");
	const first_name = claimsParser.getString("given_name");
	const middle_name = claimsParser.getString("middle_name");
	const picture = claimsParser.getString("picture");
	const email = claimsParser.getString("email");

	const existingUser = getUserFromGoogleId(googleId);
	if (existingUser !== null) {
		const sessionToken = generateSessionToken();
		const session = await createSession(sessionToken, existingUser.id, event.platform.env.GH_SESSIONS);
		setSessionTokenCookie(event, sessionToken, session.expiresAt);
		return new Response(null, {
			status: 302,
			headers: {
				Location: "/"
			}
		});
	}

	const dbSession = driverInstance.session()
	const user = createUser(dbSession, googleId, email, first_name, family_name, middle_name);
	dbSession.close();
	const sessionToken = generateSessionToken();
	const session = await createSession(sessionToken, user.id, event.platform.env.GH_SESSIONS);
	setSessionTokenCookie(event, sessionToken, session.expiresAt);

	return new Response(null, {
		status: 302,
		headers: {
			Location: "/"
		}
	});
}