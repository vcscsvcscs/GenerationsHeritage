<script lang="ts">
	import { child, date, description, edit, file, from_time, media_title, notes, parent, relation_type, sibling, spouse, title, until, upload } from '$lib/paraglide/messages';
	import type { Edge } from '@xyflow/svelte';
	import ModalButtons from '$lib/relationship/ModalButtons.svelte';
	import type { components, operations } from '$lib/api/api.gen';

	let {
		closeModal,
		onCreation = (newEdges: Edge[]) => {},
		editorMode = false,
		createRelationship = false,
		startNode = undefined,
		endNode = undefined
	} = $props<{
		closeModal: () => void;
		onCreation?: (newEdges: Edge[]) => void;
		editorMode?: boolean;
		createRelationship?: boolean;
		startNode?: string;
		endNode?: string;
	}>();

	let relationships: components['schemas']['dbtypeRelationship'][] = $state([]);
	let newRelationship: components['schemas']['FamilyRelationship'] = $state({
		verified: false,
		notes: '',
		from: '',
		to: ''
	});
	let relationshiptype: 'sibling' | 'child' | 'parent' | 'spouse' | undefined = $state('sibling');

	async function getRelationships(startId: string, endId :string) {
		if (startId === undefined || endId === undefined) {
			alert('');
			return;
		}

		const response = await fetch(`/api/relationship/${startId}/${endId}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
			console.log('Cannot get relationships, status: ' + response.status);
			return;
		}

		relationships.push((await response.json()) as components['schemas']['dbtypeRelationship']);
	}

	getRelationships(startNode,endNode);
	getRelationships(startNode,endNode);

	async function save() {
		for (const r of relationships) {
			const patchBody: components['schemas']['FamilyRelationship'] = r.Props ?? {};

			const response = await fetch(`/api/relationship/${r.StartId}/${r.EndId}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(patchBody)
			});

			if (!response.ok) {
				console.log(`Failed to save relationship ${r.StartId} → ${r.EndId}`);
			}
		}

		closeModal();
	}

	async function createNewRelationship() {
		if (relationships.length > 0) {
			alert('Relationship already exists');
			createRelationship = false;

			return;
		}

		if (!startNode || !endNode) {
			alert('Please select nodes');
			return;
		}

		let body: operations['createRelationship']['requestBody']['content']['application/json'] = {
			id1: startNode,
			id2: endNode,
			type: relationshiptype,
			relationship: newRelationship
		};

		const response = await fetch(`/api/relationship`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		});

		if (!response.ok) {
			console.log('Cannot create relationship');
			return;
		}

		const created = await response.json() as components['schemas']['dbtypeRelationship'][];
		relationships.push(...created);

		let newEdges: Edge[] = [];
		for (const r of created) {
			newEdges.push({
				id: r.ElementId!,
				source: r.StartElementId!,
				target: r.EndElementId!,
				type: 'relationship',
				data: {...r.Props, type: r.Type},
			});
		}
		onCreation(newEdges);
	}
</script>


<div class="modal modal-open z-8">
	<div class="modal-box w-full max-w-xl">
		<div class="bg-base-100 sticky top-0 z-7">
			<ModalButtons
				{editorMode}
				onCreate={createNewRelationship}
				onClose={closeModal}
				onSave={save}
				onToggleEdit={() => {
					editorMode = !editorMode;
				}}
			/>
			<div class="divider"></div>
		</div>
		{#if createRelationship}
			<!-- Relationship type selector -->
			<div class="form-control mt-4">
				<label for="relationshiptype" class="label">{relation_type()}</label>
				<select id="relationshiptype" bind:value={relationshiptype} class="select select-bordered">
					<option value="sibling">{sibling()}</option>
					<option value="child">{child()}</option>
					<option value="parent">{parent()}</option>
					<option value="spouse">{spouse()}</option>
				</select>
			</div>
		{/if}
		{#if !createRelationship}
			<!-- Editor mode: show all existing relationships -->
			{#each relationships as r, index}
				<div class="border-base-300 mt-4 rounded border p-4">
					<div class="form-control">
						{#if editorMode}
							<label for={`relationshiptype-${index}`} class="label">{relation_type()}</label>
							<select id={`relationshiptype-${index}`} bind:value={r.Type} class="select select-bordered">
								<option value="sibling">{sibling()}</option>
								<option value="child">{child()}</option>
								<option value="parent">{parent()}</option>
								<option value="spouse">{spouse()}</option>
							</select>
						{:else}
							<p><strong>{relation_type()}:</strong> {r.Type}</p>
						{/if}
					</div>
					<div class="form-control mt-2">
						{#if editorMode}
						<label for={`verified-${index}`} class="label">Verified</label>
						<input
							id={`verified-${index}`}
							type="checkbox"
							bind:value={r.Props!.verified}
							class="checkbox"
						/>
						{:else}
							<p><strong>Verified:</strong>{r.Props?.verified}</p>
						{/if}
					</div>
					<div class="form-control mt-2">ú
						{#if editorMode}
						<label for={`notes-${index}`} class="label">{notes()}</label>
						<textarea
							id={`notes-${index}`}
							bind:value={r.Props!.notes}
							class="textarea textarea-bordered w-full"
						>
						</textarea>
						{:else}
							<p><strong>{notes()}:</strong> {r.Props?.notes}</p>
						{/if}
					</div>
					<div class="form-control mt-2">
						{#if editorMode}
						<label for={`from-${index}`} class="label">{from_time()}</label>
						<input
							id={`from-${index}`}
							type="date"
							bind:value={r.Props!.from}
							class="input input-bordered w-full"
						/>
						{:else}
							<p><strong>{from_time()}:</strong> {r.Props?.from}</p>
						{/if}
					</div>
					<div class="form-control mt-2">
						{#if editorMode}
						<label for={`to-${index}`} class="label">{until()}</label>
						<input
							id={`to-${index}`}
							type="date"
							bind:value={r.Props!.to}
							class="input input-bordered w-full"
						/>
						{:else}
							<p><strong>{until()}:</strong> {r.Props?.to}</p>
						{/if}
					</div>
				</div>
			{/each}
		{:else}
			<!-- Creator mode: only one relationship -->
			<div class="border-base-300 mt-4 rounded border p-4">
				<div class="form-control">
					<label for="verified" class="label">Verified</label>
					<input
						id="verified"
						type="checkbox"
						bind:checked={newRelationship.verified}
						class="checkbox"
					/>
				</div>
				<div class="form-control mt-2">
					<label for="notes" class="label">{notes()}</label>
					<textarea
						id="notes"
						bind:value={newRelationship.notes}
						class="textarea textarea-bordered w-full"
					></textarea>
				</div>
				<div class="form-control mt-2">
					<label for="from" class="label">{from_time()}</label>
					<input
						id="from"
						type="date"
						bind:value={newRelationship.from}
						class="input input-bordered w-full"
					/>
				</div>
				<div class="form-control mt-2">
					<label for="to" class="label">{until()}</label>
					<input
						id="to"
						type="date"
						bind:value={newRelationship.to}
						class="input input-bordered w-full"
					/>
				</div>
			</div>
		{/if}
	</div>
</div>
