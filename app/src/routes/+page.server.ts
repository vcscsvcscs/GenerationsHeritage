import { fail, redirect } from "@sveltejs/kit";
import { deleteSessionTokenCookie, invalidateSession } from "$lib/server/session";
import type { Actions, RequestEvent } from "./$types";

export async function load(event: RequestEvent) {
	if (event.locals.session === null /*|| event.locals.familytree === nul*/) {
		return redirect(302, "/login");
	}
	return {
		// TODO - Add Family Graph
	};
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
		return fail(500, { message: "Server configuration error" });
	}

	deleteSessionTokenCookie(event);

	return redirect(302, "/login");
}