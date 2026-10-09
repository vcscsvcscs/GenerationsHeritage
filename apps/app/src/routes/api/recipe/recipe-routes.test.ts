import { describe, it, expect, vi, beforeEach } from 'vitest';
import { makeEvent, upstreamFail, upstreamOk } from '$lib/server/test-event';

const { client } = vi.hoisted(() => ({
	client: {
		GET: vi.fn(),
		POST: vi.fn(),
		PATCH: vi.fn(),
		DELETE: vi.fn()
	}
}));
vi.mock('$lib/api/client', () => ({ client }));

import * as recipe from './[ID]/+server';
import * as relationship from './[ID]/relationship/+server';
import * as comment from './[ID]/comment/+server';
import * as variation from './[ID]/variation/+server';
import * as variations from './[ID]/variations/+server';
import * as person from './person/[ID]/+server';
import * as cookbook from '../cookbook/+server';

type Handler = (event: never) => Promise<Response>;
const call = (handler: unknown, opts?: Parameters<typeof makeEvent>[0]) =>
	(handler as Handler)(makeEvent<never>(opts));

const headers = (userId: number) => ({ 'X-User-ID': userId });

beforeEach(() => {
	Object.values(client).forEach((m) => m.mockReset());
});

describe('authentication', () => {
	it.each([
		['recipe GET', recipe.GET],
		['recipe PATCH', recipe.PATCH],
		['recipe DELETE', recipe.DELETE],
		['relationship POST', relationship.POST],
		['relationship DELETE', relationship.DELETE],
		['comment GET', comment.GET],
		['comment POST', comment.POST],
		['comment PATCH', comment.PATCH],
		['comment DELETE', comment.DELETE],
		['variation POST', variation.POST],
		['variations GET', variations.GET],
		['person GET', person.GET],
		['person POST', person.POST],
		['cookbook GET', cookbook.GET]
	])('%s answers 401 json without a session', async (_name, handler) => {
		const res = await call(handler, { userId: null, params: { ID: '1' } });
		expect(res.status).toBe(401);
		expect(Object.values(client).every((m) => m.mock.calls.length === 0)).toBe(true);
	});
});

describe('GET /api/recipe/[ID]', () => {
	it('forwards the details payload including can_edit', async () => {
		const payload = { recipe: { Id: 5, Props: { name: 'Soup' } }, can_edit: true };
		client.GET.mockResolvedValue(upstreamOk(payload));
		const res = await call(recipe.GET, { params: { ID: '5' } });
		expect(client.GET).toHaveBeenCalledWith('/recipe/{id}', {
			params: { path: { id: 5 }, header: headers(7) }
		});
		expect(await res.json()).toEqual(payload);
	});

	it('forwards 404', async () => {
		client.GET.mockResolvedValue(upstreamFail({ msg: 'gone' }, 404));
		const res = await call(recipe.GET, { params: { ID: '5' } });
		expect(res.status).toBe(404);
		expect(await res.json()).toEqual({ msg: 'gone' });
	});
});

describe('PATCH/DELETE /api/recipe/[ID]', () => {
	it('rejects an invalid json body without calling upstream', async () => {
		const res = await call(recipe.PATCH, { params: { ID: '5' }, rawBody: '{' });
		expect(res.status).toBe(400);
		expect(client.PATCH).not.toHaveBeenCalled();
	});

	it('patches with the request body', async () => {
		client.PATCH.mockResolvedValue(upstreamOk({ Id: 5 }));
		await call(recipe.PATCH, { params: { ID: '5' }, body: { name: 'New' } });
		expect(client.PATCH).toHaveBeenCalledWith('/recipe/{id}', {
			params: { path: { id: 5 }, header: headers(7) },
			body: { name: 'New' }
		});
	});

	it('passes 404 from an already deleted recipe through', async () => {
		client.DELETE.mockResolvedValue(upstreamFail({ msg: 'not found' }, 404));
		const res = await call(recipe.DELETE, { params: { ID: '5' } });
		expect(res.status).toBe(404);
	});

	it('deletes and returns 200', async () => {
		client.DELETE.mockResolvedValue(upstreamOk({ description: 'ok' }));
		const res = await call(recipe.DELETE, { params: { ID: '5' } });
		expect(res.status).toBe(200);
	});
});

describe('/api/recipe/[ID]/relationship', () => {
	it('passes person_id and relationship through without adding an id', async () => {
		client.POST.mockResolvedValue(upstreamOk({ Id: 1 }));
		await call(relationship.POST, {
			params: { ID: '5' },
			body: { relationship: { like_it: true } }
		});
		expect(client.POST).toHaveBeenCalledWith('/recipe/{id}/relationship', {
			params: { path: { id: 5 }, header: headers(7) },
			body: { relationship: { like_it: true } }
		});
	});

	it('resolves personId=me to the session user on delete', async () => {
		client.DELETE.mockResolvedValue(upstreamOk());
		const res = await call(relationship.DELETE, {
			params: { ID: '5' },
			url: 'http://localhost/?personId=me'
		});
		expect(res.status).toBe(200);
		expect(client.DELETE).toHaveBeenCalledWith('/recipe/{id}/relationship', {
			params: { path: { id: 5 }, query: { personId: 7 }, header: headers(7) }
		});
	});

	it('passes an explicit personId and rejects a missing or invalid one', async () => {
		client.DELETE.mockResolvedValue(upstreamOk());
		await call(relationship.DELETE, { params: { ID: '5' }, url: 'http://localhost/?personId=9' });
		expect(client.DELETE.mock.calls[0][1].params.query).toEqual({ personId: 9 });

		for (const url of ['http://localhost/', 'http://localhost/?personId=abc']) {
			const res = await call(relationship.DELETE, { params: { ID: '5' }, url });
			expect(res.status).toBe(400);
		}
		expect(client.DELETE).toHaveBeenCalledTimes(1);
	});
});

describe('/api/recipe/[ID]/comment', () => {
	const flat = {
		comment: { message: 'yum', sent_at: 1_700_000_000, edited: null },
		commenter: { id: 7, first_name: 'A' }
	};

	it('GET forwards the flat comments list', async () => {
		client.GET.mockResolvedValue(upstreamOk({ comments: [flat] }));
		const res = await call(comment.GET, { params: { ID: '5' } });
		expect(client.GET).toHaveBeenCalledWith('/recipe/{id}/comment', {
			params: { path: { id: 5 }, header: headers(7) }
		});
		expect(await res.json()).toEqual({ comments: [flat] });
	});

	it('POST sends only the message (no timestamps from the client)', async () => {
		client.POST.mockResolvedValue(upstreamOk(flat));
		const res = await call(comment.POST, {
			params: { ID: '5' },
			body: { message: 'yum', sent_at: 'x', edited: 'y' }
		});
		expect(client.POST).toHaveBeenCalledWith('/recipe/{id}/comment', {
			params: { path: { id: 5 }, header: headers(7) },
			body: { message: 'yum' }
		});
		expect(await res.json()).toEqual(flat);
	});

	it.each([{}, { message: '   ' }, { message: 5 }])('POST rejects %j with 400', async (body) => {
		const res = await call(comment.POST, { params: { ID: '5' }, body });
		expect(res.status).toBe(400);
		expect(client.POST).not.toHaveBeenCalled();
	});

	it('PATCH returns the full comment and commenter', async () => {
		const edited = { ...flat, comment: { ...flat.comment, edited: 1_700_000_100 } };
		client.PATCH.mockResolvedValue(upstreamOk(edited));
		const res = await call(comment.PATCH, { params: { ID: '5' }, body: { message: 'yummier' } });
		expect(client.PATCH).toHaveBeenCalledWith('/recipe/{id}/comment', {
			params: { path: { id: 5 }, header: headers(7) },
			body: { message: 'yummier' }
		});
		expect(await res.json()).toEqual(edited);
	});

	it('PATCH forwards 404 when the user has no comment', async () => {
		client.PATCH.mockResolvedValue(upstreamFail({ msg: 'none' }, 404));
		const res = await call(comment.PATCH, { params: { ID: '5' }, body: { message: 'x' } });
		expect(res.status).toBe(404);
	});

	it('DELETE forwards to upstream', async () => {
		client.DELETE.mockResolvedValue(upstreamOk({ description: 'deleted' }));
		const res = await call(comment.DELETE, { params: { ID: '5' } });
		expect(client.DELETE).toHaveBeenCalledWith('/recipe/{id}/comment', {
			params: { path: { id: 5 }, header: headers(7) }
		});
		expect(res.status).toBe(200);
	});
});

describe('variations', () => {
	it('GET forwards the variations list', async () => {
		const payload = {
			variations: [
				{
					variation: { Id: 8 },
					variation_relationship: { notes: 'n', created_at: 1 },
					creator: { id: 2 }
				}
			]
		};
		client.GET.mockResolvedValue(upstreamOk(payload));
		const res = await call(variations.GET, { params: { ID: '5' } });
		expect(client.GET).toHaveBeenCalledWith('/recipe/{id}/variations', {
			params: { path: { id: 5 }, header: headers(7) }
		});
		expect(await res.json()).toEqual(payload);
	});

	it('POST forwards recipe and notes', async () => {
		const body = { recipe: { name: 'v' }, variation_notes: 'less salt' };
		client.POST.mockResolvedValue(upstreamOk({ recipe: { Id: 9 } }));
		await call(variation.POST, { params: { ID: '5' }, body });
		expect(client.POST).toHaveBeenCalledWith('/recipe/{id}/variation', {
			params: { path: { id: 5 }, header: headers(7) },
			body
		});
	});

	it('POST rejects a missing recipe', async () => {
		const res = await call(variation.POST, { params: { ID: '5' }, body: {} });
		expect(res.status).toBe(400);
		expect(client.POST).not.toHaveBeenCalled();
	});
});

describe('/api/recipe/person/[ID]', () => {
	it('GET resolves me and forwards the entries payload', async () => {
		const payload = {
			entries: [{ recipe: { Id: 1 }, relationship: null, created: true, can_edit: true }]
		};
		client.GET.mockResolvedValue(upstreamOk(payload));
		const res = await call(person.GET, { params: { ID: 'me' } });
		expect(client.GET).toHaveBeenCalledWith('/person/{id}/recipes', {
			params: { path: { id: 7 }, header: headers(7) }
		});
		expect(await res.json()).toEqual(payload);
	});

	it('GET rejects a non numeric person', async () => {
		const res = await call(person.GET, { params: { ID: 'abc' } });
		expect(res.status).toBe(400);
	});

	it('POST forwards recipe and relationship', async () => {
		const body = { recipe: { name: 'r' }, relationship: { favourite: true } };
		client.POST.mockResolvedValue(upstreamOk({}));
		await call(person.POST, { params: { ID: '3' }, body });
		expect(client.POST).toHaveBeenCalledWith('/person/{id}/recipes', {
			params: { path: { id: 3 }, header: headers(7) },
			body
		});
	});
});

describe('/api/cookbook', () => {
	it('defaults the distance to 3 and forwards entries', async () => {
		const payload = { entries: [{ recipe: { Id: 1 }, relationship: null, can_edit: false }] };
		client.GET.mockResolvedValue(upstreamOk(payload));
		const res = await call(cookbook.GET);
		expect(client.GET).toHaveBeenCalledWith('/cookbook', {
			params: { query: { distance: 3 }, header: headers(7) }
		});
		expect(await res.json()).toEqual(payload);
	});

	it('rejects an invalid distance', async () => {
		const res = await call(cookbook.GET, { url: 'http://localhost/?distance=abc' });
		expect(res.status).toBe(400);
	});
});
