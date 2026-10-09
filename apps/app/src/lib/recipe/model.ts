import type { components } from '$lib/api/api.gen';

type Schemas = components['schemas'];

export type Recipe = Schemas['Recipe'];
export type RecipeProperties = Schemas['RecipeProperties'];
export type RecipeEntry = Schemas['RecipeEntry'];
export type CookbookEntry = Schemas['CookbookEntry'];
export type RecipeComment = Schemas['RecipeComment'];
export type RecipeVariation = Schemas['RecipeVariation'];

export type RecipeDraft = {
	name: string;
	origin: string;
	category: string;
	description: string;
	ingredients: string[];
	instructions: string[];
	notes: string;
};

export const emptyDraft = (): RecipeDraft => ({
	name: '',
	origin: '',
	category: '',
	description: '',
	ingredients: [''],
	instructions: [''],
	notes: ''
});

export const toDraft = (props: RecipeProperties | undefined | null): RecipeDraft => ({
	name: props?.name ?? '',
	origin: props?.origin ?? '',
	category: props?.category ?? '',
	description: props?.description ?? '',
	ingredients: [...(props?.ingredients ?? [])],
	instructions: [...(props?.instructions ?? [])],
	notes: props?.notes ?? ''
});

export const cleanDraft = (draft: RecipeDraft): RecipeDraft => ({
	...draft,
	ingredients: draft.ingredients.filter((i) => i.trim() !== ''),
	instructions: draft.instructions.filter((i) => i.trim() !== '')
});

type NamedPerson = {
	first_name?: string | null;
	last_name?: string | null;
};

export const personName = (person: NamedPerson | null | undefined): string =>
	[person?.first_name, person?.last_name].filter(Boolean).join(' ') || '?';

export const isLiked = (entry: Pick<RecipeEntry, 'relationship'> | null | undefined): boolean =>
	entry?.relationship != null;

export const likedRecipeIds = (entries: RecipeEntry[] | null | undefined): Set<number> =>
	new Set(
		(entries ?? [])
			.filter((entry) => isLiked(entry) && entry.recipe?.Id !== undefined)
			.map((entry) => entry.recipe!.Id!)
	);

export const canEdit = (entry: { can_edit?: boolean } | null | undefined): boolean =>
	entry?.can_edit === true;

export const formatUnixSeconds = (
	seconds: number | null | undefined,
	locale?: string
): string | null =>
	seconds === null || seconds === undefined || !Number.isFinite(seconds)
		? null
		: new Date(seconds * 1000).toLocaleString(locale);

export const ownComment = (
	comments: RecipeComment[],
	userId: number | null | undefined
): RecipeComment | undefined =>
	userId === null || userId === undefined
		? undefined
		: comments.find((c) => c.commenter?.id === userId);

export const mergeComment = (comments: RecipeComment[], updated: RecipeComment): RecipeComment[] =>
	comments.some((c) => c.commenter?.id === updated.commenter?.id)
		? comments.map((c) => (c.commenter?.id === updated.commenter?.id ? updated : c))
		: [...comments, updated];
