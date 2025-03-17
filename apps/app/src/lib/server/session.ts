import type { KVNamespace } from "@cloudflare/workers-types";
import { encodeBase32, encodeHexLowerCase } from "@oslojs/encoding";
import { sha256 } from "@oslojs/crypto/sha2";

import type { RequestEvent } from "@sveltejs/kit";

// in seconds
const EXPIRATION_TTL: number = 60 * 60 * 24 * 7;

export async function validateSessionToken(token: string, sessions: KVNamespace): Promise<SessionValidationResult> {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
	const session: Session | null = await sessions.get(sessionId, { type: "json" });

	if (!session) {
		return null;
	}

	if (Date.now() >= session.expiresAt.getTime() - 1000 * 60 * 60 * 24 * 15) {
		await sessions.put(sessionId, JSON.stringify(session), { expirationTtl: EXPIRATION_TTL });
	}
	return session;
}

export async function invalidateSession(sessionId: string, sessions: KVNamespace): Promise<void> {
	await sessions.delete(sessionId);
}

export async function invalidateUserSessions(userId: number, sessions: KVNamespace): Promise<void> {
	const keys = await sessions.list({ prefix: `${userId}:` });
	for (const key of keys.keys) {
		await sessions.delete(key.name);
	}
}

export function setSessionTokenCookie(event: RequestEvent, token: string, expiresAt: Date): void {
	event.cookies.set("session", token, {
		httpOnly: true,
		path: "/",
		secure: import.meta.env.PROD,
		sameSite: "lax",
		expires: expiresAt
	});
}

export function deleteSessionTokenCookie(event: RequestEvent): void {
	event.cookies.set("session", "", {
		httpOnly: true,
		path: "/",
		secure: import.meta.env.PROD,
		sameSite: "lax",
		maxAge: 0
	});
}

export function generateSessionToken(): string {
	const tokenBytes = new Uint8Array(20);
	crypto.getRandomValues(tokenBytes);
	const token = encodeBase32(tokenBytes).toLowerCase();
	return token;
}

export async function createSession(token: string, userId: string, sessions: KVNamespace): Promise<Session> {
	const sessionId = `${userId}:${encodeHexLowerCase(sha256(new TextEncoder().encode(token)))}`;
	const session: Session = {
		id: sessionId,
		userId,
		expiresAt: new Date(Date.now() + 1000 * EXPIRATION_TTL)
	};
	await sessions.put(sessionId, JSON.stringify(session), { expirationTtl: EXPIRATION_TTL });
	return session;
}

export interface Session {
	id: string;
	expiresAt: Date;
	userId: string;
}

type SessionValidationResult = Session | null;
