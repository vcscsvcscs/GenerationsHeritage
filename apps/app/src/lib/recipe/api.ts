import type { components } from '$lib/api/api.gen';
import type { RecipeProperties } from './model';

type Schemas = components['schemas'];

export type ApiResult<T> = { ok: true; data: T } | { ok: false; status: number };

async function call<T>(url: string, init?: RequestInit): Promise<ApiResult<T>> {
	try {
		const response = await fetch(url, init);
		if (!response.ok) return { ok: false, status: response.status };
		const text = await response.text();
		return { ok: true, data: (text ? JSON.parse(text) : undefined) as T };
	} catch (e) {
		console.error('Recipe API request failed:', url, e);
		return { ok: false, status: 0 };
	}
}

const json = (method: string, body: unknown): RequestInit => ({
	method,
	headers: { 'Content-Type': 'application/json' },
	body: JSON.stringify(body)
});

export const getPersonRecipes = (personId: number | 'me') =>
	call<Schemas['PersonRecipes']>(`/api/recipe/person/${personId}`);

export const getCookbook = (distance: number) =>
	call<Schemas['Cookbook']>(`/api/cookbook?distance=${distance}`);

export const getRecipe = (id: number) => call<Schemas['RecipeDetails']>(`/api/recipe/${id}`);

export const createRecipe = (
	personId: number,
	recipe: RecipeProperties,
	relationship?: Schemas['LikesProperties']
) => call<unknown>(`/api/recipe/person/${personId}`, json('POST', { recipe, relationship }));

export const updateRecipe = (id: number, recipe: RecipeProperties) =>
	call<unknown>(`/api/recipe/${id}`, json('PATCH', recipe));

export const deleteRecipe = (id: number) =>
	call<unknown>(`/api/recipe/${id}`, { method: 'DELETE' });

export const likeRecipe = (id: number) =>
	call<Schemas['Likes']>(
		`/api/recipe/${id}/relationship`,
		json('POST', { relationship: { like_it: true } })
	);

export const unlikeRecipe = (id: number) =>
	call<unknown>(`/api/recipe/${id}/relationship?personId=me`, { method: 'DELETE' });

export const getComments = (recipeId: number) =>
	call<Schemas['RecipeComments']>(`/api/recipe/${recipeId}/comment`);

export const postComment = (recipeId: number, message: string) =>
	call<Schemas['RecipeComment']>(`/api/recipe/${recipeId}/comment`, json('POST', { message }));

export const editComment = (recipeId: number, message: string) =>
	call<Schemas['RecipeComment']>(`/api/recipe/${recipeId}/comment`, json('PATCH', { message }));

export const deleteComment = (recipeId: number) =>
	call<unknown>(`/api/recipe/${recipeId}/comment`, { method: 'DELETE' });

export const getVariations = (recipeId: number) =>
	call<{ variations?: Schemas['RecipeVariation'][] }>(`/api/recipe/${recipeId}/variations`);

export const createVariation = (recipeId: number, recipe: RecipeProperties, notes: string) =>
	call<{ recipe?: Schemas['Recipe'] }>(
		`/api/recipe/${recipeId}/variation`,
		json('POST', { recipe, variation_notes: notes.trim() === '' ? null : notes.trim() })
	);
