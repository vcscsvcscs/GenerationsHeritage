<script lang="ts">
	import { onMount } from 'svelte';
	import {
		comments as commentsLabel,
		add_comment,
		post_comment,
		no_comments,
		comment_edited,
		delete_action as deleteLabel,
		edit,
		save,
		cancel
	} from '$lib/paraglide/messages';
	import { languageTag } from '$lib/paraglide/runtime';
	import { getComments, postComment, editComment, deleteComment } from './api';
	import { errorMessage } from './errors';
	import { formatUnixSeconds, mergeComment, ownComment, personName } from './model';
	import type { RecipeComment } from './model';

	let { recipeId, currentUserId }: { recipeId: number; currentUserId: number | null } = $props();

	let comments: RecipeComment[] = $state([]);
	let isLoading = $state(true);
	let busy = $state(false);
	let error: string | null = $state(null);
	let newMessage = $state('');
	let editing = $state(false);
	let editMessage = $state('');

	const mine = $derived(ownComment(comments, currentUserId));

	onMount(async () => {
		const result = await getComments(recipeId);
		if (result.ok) {
			comments = result.data.comments ?? [];
		} else {
			error = errorMessage(result.status);
		}
		isLoading = false;
	});

	async function submit() {
		if (busy || newMessage.trim() === '') return;
		busy = true;
		error = null;
		const result = await postComment(recipeId, newMessage);
		if (result.ok) {
			comments = mergeComment(comments, result.data);
			newMessage = '';
		} else {
			error = errorMessage(result.status);
		}
		busy = false;
	}

	async function saveEdit() {
		if (busy || editMessage.trim() === '') return;
		busy = true;
		error = null;
		const result = await editComment(recipeId, editMessage);
		if (result.ok) {
			comments = mergeComment(comments, result.data);
			editing = false;
		} else {
			error = errorMessage(result.status);
		}
		busy = false;
	}

	async function remove() {
		if (busy) return;
		busy = true;
		error = null;
		const result = await deleteComment(recipeId);
		if (result.ok || result.status === 404) {
			comments = comments.filter((c) => c.commenter?.id !== currentUserId);
			editing = false;
		} else {
			error = errorMessage(result.status);
		}
		busy = false;
	}
</script>

<section class="flex flex-col gap-2">
	<h4 class="text-base font-bold">{commentsLabel()}</h4>

	{#if error}
		<div class="alert alert-error alert-soft py-2 text-sm" role="alert">{error}</div>
	{/if}

	{#if isLoading}
		<span class="loading loading-spinner"></span>
	{:else if comments.length === 0}
		<p class="text-base-content/60 text-sm">{no_comments()}</p>
	{:else}
		<ul class="flex flex-col gap-2">
			{#each comments as c (c.commenter?.id)}
				<li class="bg-base-200 rounded-box p-3">
					<div class="text-base-content/60 flex items-center gap-2 text-xs">
						{#if c.commenter?.profile_picture}
							<img src={c.commenter.profile_picture} alt="" class="h-5 w-5 rounded-full" />
						{/if}
						<span class="font-semibold">{personName(c.commenter)}</span>
						<span>{formatUnixSeconds(c.comment?.sent_at, languageTag()) ?? ''}</span>
						{#if c.comment?.edited}
							<span class="italic">
								({comment_edited()}
								{formatUnixSeconds(c.comment.edited, languageTag())})
							</span>
						{/if}
					</div>
					{#if editing && c.commenter?.id === currentUserId}
						<textarea
							class="textarea textarea-bordered mt-2 w-full"
							rows="2"
							bind:value={editMessage}
						></textarea>
						<div class="mt-1 flex justify-end gap-2">
							<button class="btn btn-ghost btn-xs" onclick={() => (editing = false)}>
								{cancel()}
							</button>
							<button
								class="btn btn-primary btn-xs"
								onclick={saveEdit}
								disabled={busy || editMessage.trim() === ''}
							>
								{save()}
							</button>
						</div>
					{:else}
						<p class="mt-1 text-sm whitespace-pre-wrap">{c.comment?.message ?? ''}</p>
						{#if c.commenter?.id === currentUserId}
							<div class="mt-1 flex justify-end gap-2">
								<button
									class="btn btn-secondary btn-xs"
									onclick={() => {
										editMessage = c.comment?.message ?? '';
										editing = true;
									}}
									disabled={busy}
								>
									{edit()}
								</button>
								<button class="btn btn-error btn-xs" onclick={remove} disabled={busy}>
									{deleteLabel()}
								</button>
							</div>
						{/if}
					{/if}
				</li>
			{/each}
		</ul>
	{/if}

	{#if !isLoading && !mine && currentUserId !== null}
		<textarea
			class="textarea textarea-bordered w-full"
			rows="2"
			placeholder={add_comment()}
			bind:value={newMessage}
		></textarea>
		<div class="flex justify-end">
			<button
				class="btn btn-primary btn-sm"
				onclick={submit}
				disabled={busy || newMessage.trim() === ''}
			>
				{post_comment()}
			</button>
		</div>
	{/if}
</section>
