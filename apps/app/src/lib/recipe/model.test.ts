import { describe, it, expect } from 'vitest';
import {
	canEdit,
	cleanDraft,
	formatUnixSeconds,
	isLiked,
	likedRecipeIds,
	mergeComment,
	ownComment,
	personName,
	toDraft
} from './model';
import type { RecipeComment, RecipeEntry } from './model';

describe('personName', () => {
	it('joins first and last name', () => {
		expect(personName({ first_name: 'Anna', last_name: 'Kovacs' })).toBe('Anna Kovacs');
	});
	it.each([undefined, null, {}, { first_name: null, last_name: '' }])(
		'falls back to ? for %j',
		(p) => {
			expect(personName(p)).toBe('?');
		}
	);
});

describe('likes', () => {
	const entries: RecipeEntry[] = [
		{ recipe: { Id: 1 }, relationship: { Props: { like_it: true } }, created: false },
		{ recipe: { Id: 2 }, relationship: null, created: true },
		{ recipe: { Id: 3 } },
		{ relationship: {} }
	];

	it('treats only a non-null relationship as liked', () => {
		expect(entries.map(isLiked)).toEqual([true, false, false, true]);
	});

	it('collects liked recipe ids, ignoring entries without a recipe', () => {
		expect(likedRecipeIds(entries)).toEqual(new Set([1]));
		expect(likedRecipeIds(undefined)).toEqual(new Set());
	});
});

describe('canEdit', () => {
	it('is true only when can_edit is explicitly true', () => {
		expect(canEdit({ can_edit: true })).toBe(true);
		expect(canEdit({ can_edit: false })).toBe(false);
		expect(canEdit({})).toBe(false);
		expect(canEdit(undefined)).toBe(false);
	});
});

describe('draft helpers', () => {
	it('normalises null and missing properties', () => {
		expect(toDraft({ name: 'Soup', origin: null, ingredients: null })).toEqual({
			name: 'Soup',
			origin: '',
			category: '',
			description: '',
			ingredients: [],
			instructions: [],
			notes: ''
		});
		expect(toDraft(null).name).toBe('');
	});

	it('copies arrays so editing the draft does not mutate the source', () => {
		const props = { ingredients: ['a'] };
		toDraft(props).ingredients.push('b');
		expect(props.ingredients).toEqual(['a']);
	});

	it('drops blank ingredients and instructions', () => {
		const cleaned = cleanDraft({
			...toDraft({ name: 'x' }),
			ingredients: ['a', ' ', ''],
			instructions: ['', 'step']
		});
		expect(cleaned.ingredients).toEqual(['a']);
		expect(cleaned.instructions).toEqual(['step']);
	});
});

describe('formatUnixSeconds', () => {
	it('interprets the value as seconds', () => {
		expect(formatUnixSeconds(1_700_000_000, 'en-US')).toBe(
			new Date(1_700_000_000_000).toLocaleString('en-US')
		);
	});
	it.each([null, undefined, NaN])('returns null for %s', (v) => {
		expect(formatUnixSeconds(v)).toBeNull();
	});
});

describe('comments', () => {
	const mine: RecipeComment = { comment: { message: 'a' }, commenter: { id: 1 } };
	const other: RecipeComment = { comment: { message: 'b' }, commenter: { id: 2 } };

	it('finds the comment of the current user', () => {
		expect(ownComment([other, mine], 1)).toBe(mine);
		expect(ownComment([other], 1)).toBeUndefined();
		expect(ownComment([mine], null)).toBeUndefined();
	});

	it('replaces an existing comment by commenter or appends a new one', () => {
		const edited = { comment: { message: 'a2', edited: 5 }, commenter: { id: 1 } };
		expect(mergeComment([mine, other], edited)).toEqual([edited, other]);
		expect(mergeComment([other], mine)).toEqual([other, mine]);
	});
});
