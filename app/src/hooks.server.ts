import type { Handle } from '@sveltejs/kit';
import { i18n } from '$lib/i18n';
import { validateSessionToken, setSessionTokenCookie, deleteSessionTokenCookie } from "$lib/server/session";
import { sequence } from "@sveltejs/kit/hooks";

const handleParaglide: Handle = i18n.handle();

const authHandle: Handle = async ({ event, resolve }) => {
    const token = event.cookies.get("session") ?? null;
    if (token === null) {
        event.locals.session = null;
        return resolve(event);
    }

    if(!event.platform || !event.platform.env || !event.platform.env.GH_SESSIONS){
		return new Response("Server configuration error. GH_SESSIONS KeyValue store missing", {
			status: 500
		});
	}

    const session = await validateSessionToken(token, event.platform.env.GH_SESSIONS);
    if (session !== null) {
        setSessionTokenCookie(event, token, session.expiresAt);
    } else {
        deleteSessionTokenCookie(event);
    }

    event.locals.session = session;
    return resolve(event);
};

export const handle: Handle = sequence(handleParaglide, authHandle);