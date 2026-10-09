import { describe, it, expect, vi } from 'vitest';
import { readJson, toResponse, withSession, withSessionAndId } from './proxy';

type Ev = { locals: App.Locals; params: { ID: string } };
const ev = (userId: number | null, ID = '5'): Ev => ({
	locals: { session: userId === null ? null : { id: 's', expiresAt: 0, userId } },
	params: { ID }
});

describe('toResponse', () => {
	it('serialises data on success', async () => {
		const res = toResponse({ data: { a: 1 }, response: new Response(null, { status: 200 }) });
		expect(res.status).toBe(200);
		expect(res.headers.get('Content-Type')).toBe('application/json');
		expect(await res.json()).toEqual({ a: 1 });
	});

	it('serialises the error and keeps the status on failure', async () => {
		const res = toResponse({
			error: { msg: 'nope' },
			response: new Response(null, { status: 409 })
		});
		expect(res.status).toBe(409);
		expect(await res.json()).toEqual({ msg: 'nope' });
	});

	it('returns an empty body when there is no payload', async () => {
		const res = toResponse({ response: new Response(null, { status: 200 }) });
		expect(res.status).toBe(200);
		expect(await res.text()).toBe('');
	});
});

describe('readJson', () => {
	it('returns null for invalid json', async () => {
		expect(await readJson(new Request('http://x', { method: 'POST', body: '{' }))).toBeNull();
	});
});

describe('withSession', () => {
	it('answers 401 without a session and skips the handler', async () => {
		const handler = vi.fn();
		const res = await withSession<Ev>(handler)(ev(null));
		expect(res.status).toBe(401);
		expect(handler).not.toHaveBeenCalled();
	});

	it('passes the user id to the handler', async () => {
		const handler = vi.fn().mockResolvedValue(new Response('ok'));
		await withSession<Ev>(handler)(ev(42));
		expect(handler).toHaveBeenCalledWith(expect.anything(), 42);
	});
});

describe('withSessionAndId', () => {
	it.each(['abc', '0', '-1', '1.5', ''])('rejects invalid id %j with 400', async (id) => {
		const handler = vi.fn();
		const res = await withSessionAndId<Ev>(handler)(ev(1, id));
		expect(res.status).toBe(400);
		expect(handler).not.toHaveBeenCalled();
	});

	it('passes numeric id', async () => {
		const handler = vi.fn().mockResolvedValue(new Response('ok'));
		await withSessionAndId<Ev>(handler)(ev(3, '9'));
		expect(handler).toHaveBeenCalledWith(expect.anything(), 3, 9);
	});

	it('answers 401 before validating the id', async () => {
		const res = await withSessionAndId<Ev>(vi.fn())(ev(null, 'abc'));
		expect(res.status).toBe(401);
	});
});
