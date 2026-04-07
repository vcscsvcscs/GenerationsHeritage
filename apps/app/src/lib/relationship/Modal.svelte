<script lang="ts">
	import {
		child,
		from_time,
		id,
		notes,
		parent,
		relation,
		relation_type,
		verified,
		sibling,
		spouse,
		until
	} from '$lib/paraglide/messages';
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

	async function getRelationships(startId: string, endId: string) {
		if (
			startId === undefined ||
			endId === undefined ||
			startId === '' ||
			endId === '' ||
			startId === endId
		) {
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

	if (!createRelationship) {
		getRelationships(startNode, endNode);
		getRelationships(endNode, startNode);
	}

	async function save() {
		for (const r of relationships) {
			if (!r.Props) {
				console.log('No properties found for relationship', r);
				continue;
			}
			if (r.Props.verified === undefined) {
				r.Props.verified = false;
			}
			console.debug('Saving relationship', r.StartId, r.EndId, r.Props);
			const patchBody: components['schemas']['FamilyRelationship'] = r.Props!;

			const response = await fetch(`/api/relationship/${r.StartId}/${r.EndId}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(patchBody)
			});

			if (!response.ok) {
				console.error(`Failed to save relationship ${r.StartId} → ${r.EndId}`);
			}

			if (response.status === 200) {
				console.debug(`Relationship ${r.StartId} → ${r.EndId} saved successfully`);
			} else {
				console.error(`Failed to save relationship ${r.StartId} → ${r.EndId}`);
			}
		}
		editorMode = !editorMode;
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
			id1: Number(startNode),
			id2: Number(endNode),
			type: relationshiptype,
			relationship: newRelationship
		};

		console.log('Creating relationship', body);
		const response = await fetch(`/api/relationship`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		});

		if (!response.ok) {
			console.error(
				'Cannot create relationship' + ', status: ' + response.status + (await response.json())
			);
			return;
		}

		const created = (await response.json()) as components['schemas']['dbtypeRelationship'][];
		console.debug('Relationship created successfully', created);
		relationships.push(...created);

		let newEdges: Edge[] = [];
		for (const r of created) {
			newEdges.push({
				id: r.ElementId!,
				source: r.StartElementId!,
				target: r.EndElementId!,
				type: 'relationship',
				data: { ...r.Props, type: r.Type }
			});
		}
		onCreation(newEdges);
		closeModal();
	}
</script>

<div class="modal modal-open z-8">
	<div class="modal-box w-full max-w-xl gap-4">
		<div class="bg-base-100 sticky top-0 z-7">
			<ModalButtons
				{editorMode}
				createMode={createRelationship}
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
			<div class="form-control mt-1">
				<p><strong>{id().toLowerCase()}:</strong>{startNode}</p>
			</div>
			<div class="form-control mt-1">
				<label for="endNode" class="label">{relation() + ' ' + id().toLowerCase()}:</label>
				<input id="endNode" type="text" bind:value={endNode} class="input input-bordered w-full" />
			</div>
		{/if}
		{#if !createRelationship}
			<div class="form-control mt-2">
				<p>
					<strong>{id()} 1:</strong>
					{startNode}
				</p>
				<p>
					<strong>{id()} 2:</strong>
					{endNode}
				</p>
			</div>
			<!-- Editor mode: show all existing relationships -->
			{#each relationships as r, index}
				<div class="border-base-300 mt-4 rounded border p-4">
					<div class="form-control">
						<p><strong>{relation_type()}:</strong> {r.Type}</p>
					</div>
					<div class="form-control mt-2">
						{#if editorMode}
							<label for={`verified-${index}`} class="label">Verified</label>
							<input
								id={`verified-${index}`}
								type="checkbox"
								bind:checked={relationships[index].Props!.verified}
								class="checkbox"
							/>
						{:else}
							<p><strong>{verified()}:</strong>{r.Props?.verified}</p>
						{/if}
					</div>
					<div class="form-control mt-2">
						{#if editorMode}
							<label for={`notes-${index}`} class="label">{notes()}</label>
							<textarea
								id={`notes-${index}`}
								bind:value={relationships[index].Props!.notes}
								class="textarea textarea-bordered w-full"
							></textarea>
						{:else}
							<p><strong>{notes()}:</strong> {relationships[index].Props?.notes}</p>
						{/if}
					</div>
					<div class="form-control mt-2">
						{#if editorMode}
							<label for={`from-${index}`} class="label">{from_time()}</label>
							<input
								id={`from-${index}`}
								type="date"
								bind:value={relationships[index].Props!.from}
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
								bind:value={relationships[index].Props!.to}
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
