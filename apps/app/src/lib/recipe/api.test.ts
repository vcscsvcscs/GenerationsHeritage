import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import * as api from './api';

const fetchMock = vi.fn();

beforeEach(() => {
	vi.stubGlobal('fetch', fetchMock);
	vi.spyOn(console, 'error').mockImplementation(() => {});
});
afterEach(() => {
	fetchMock.mockReset();
	vi.unstubAllGlobals();
	vi.restoreAllMocks();
});

const reply = (body: unknown, status = 200) =>
	fetchMock.mockResolvedValue(
		new Response(body === undefined ? null : JSON.stringify(body), { status })
	);

describe('result handling', () => {
	it('returns parsed data on success', async () => {
		reply({ entries: [] });
		expect(await api.getCookbook(5)).toEqual({ ok: true, data: { entries: [] } });
		expect(fetchMock).toHaveBeenCalledWith('/api/cookbook?distance=5', undefined);
	});

	it('tolerates empty success bodies', async () => {
		reply(undefined);
		expect(await api.deleteComment(1)).toEqual({ ok: true, data: undefined });
	});

	it('returns the status on http failure', async () => {
		reply({ msg: 'x' }, 401);
		expect(await api.getRecipe(1)).toEqual({ ok: false, status: 401 });
	});

	it('maps network errors to status 0', async () => {
		fetchMock.mockRejectedValue(new Error('offline'));
		expect(await api.getRecipe(1)).toEqual({ ok: false, status: 0 });
	});
});

describe('requests', () => {
	const lastCall = () => {
		const [url, init] = fetchMock.mock.calls.at(-1)!;
		return { url, method: init?.method, body: init?.body ? JSON.parse(init.body) : undefined };
	};

	it('likes with a plain LikesProperties relationship and no id', async () => {
		reply({});
		await api.likeRecipe(4);
		expect(lastCall()).toEqual({
			url: '/api/recipe/4/relationship',
			method: 'POST',
			body: { relationship: { like_it: true } }
		});
	});

	it('unlikes the current user', async () => {
		reply(undefined);
		await api.unlikeRecipe(4);
		expect(lastCall()).toMatchObject({
			url: '/api/recipe/4/relationship?personId=me',
			method: 'DELETE'
		});
	});

	it('posts and edits comments with a flat message body', async () => {
		reply({});
		await api.postComment(4, 'tasty');
		expect(lastCall()).toEqual({
			url: '/api/recipe/4/comment',
			method: 'POST',
			body: { message: 'tasty' }
		});
		await api.editComment(4, 'tastier');
		expect(lastCall()).toEqual({
			url: '/api/recipe/4/comment',
			method: 'PATCH',
			body: { message: 'tastier' }
		});
	});

	it('creates a variation with trimmed notes, or null when blank', async () => {
		reply({});
		await api.createVariation(4, { name: 'v' }, '  less salt ');
		expect(lastCall()).toEqual({
			url: '/api/recipe/4/variation',
			method: 'POST',
			body: { recipe: { name: 'v' }, variation_notes: 'less salt' }
		});
		await api.createVariation(4, { name: 'v' }, '   ');
		expect(lastCall().body.variation_notes).toBeNull();
	});

	it('creates a recipe for a person with the like properties', async () => {
		reply({});
		await api.createRecipe(9, { name: 'r' }, { favourite: true });
		expect(lastCall()).toEqual({
			url: '/api/recipe/person/9',
			method: 'POST',
			body: { recipe: { name: 'r' }, relationship: { favourite: true } }
		});
	});

	it('uses the me alias for the own recipe list', async () => {
		reply({ entries: [] });
		await api.getPersonRecipes('me');
		expect(lastCall().url).toBe('/api/recipe/person/me');
	});
});
