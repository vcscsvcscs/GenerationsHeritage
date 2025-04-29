<script lang="ts">
	import { fade } from 'svelte/transition';
	import {
		create,
		close,
		born,
		mothers_first_name,
		mothers_last_name,
		last_name,
		first_name,
		email,
		biological_sex,
		male,
		female,
		other,
		intersex,
		create_relationship_and_person,
		child,
		sibling,
		parent,
		spouse,
		relation,
		relation_type,
		notes,
		until,
		optional_field,
		from_time
	} from '$lib/paraglide/messages';
	import { onMount } from 'svelte';
	import type { components } from '$lib/api/api.gen.js';
	import { validatePersonRegistration, validateFamilyRelationship } from './validate_fields';
	import type { Node, Edge } from '@xyflow/svelte';

	let {
		closeModal = () => {},
		onCreation = (nodes: Array<Node> | null, edges: Array<Edge> | null) => {},
		onOnlyPersonCreation = (person: components['schemas']['Person']) => {},
		relationshipStartID
	}: {
		closeModal: () => void;
		onCreation: (newNodes: Array<Node> | null, newEdges: Array<Edge> | null) => void;
		onOnlyPersonCreation: (person: components['schemas']['Person']) => void | undefined;
		relationshipStartID: number | null;
	} = $props();

	let birth_date: HTMLInputElement;
	let relationship_from_time: HTMLInputElement = $state({} as HTMLInputElement);
	let relationship_until: HTMLInputElement = $state({} as HTMLInputElement);

	let draftRelationship: (components['schemas']['FamilyRelationship'] & { type: string }) | null =
		$state({} as components['schemas']['FamilyRelationship'] & { type: string });
	let draftPerson: components['schemas']['PersonRegistration'] = $state(
		{} as components['schemas']['PersonRegistration']
	);
	let error: string | undefined | null = $state();

	function onClose() {
		closeModal();
	}

	async function onCreate(event: SubmitEvent) {
		event.preventDefault();
		error = validatePersonRegistration(draftPerson);
		if (error) {
			return;
		}

		if (relationshipStartID !== null) {
			if (draftRelationship !== null) {
				error = validateFamilyRelationship(draftRelationship);
				if (error) {
					return;
				}
			}

			let requestBody = {
				relationship: draftRelationship,
				type: draftRelationship!.type,
				person: draftPerson
			} as {
				person: components['schemas']['PersonRegistration'];
				type?: 'child' | 'parent' | 'spouse' | 'sibling';
				relationship: components['schemas']['FamilyRelationship'];
			};
			let response = await fetch(`/api/person_and_relationship/${relationshipStartID}`, {
				method: 'POST',
				body: JSON.stringify(requestBody),
				headers: {
					'Content-Type': 'application/json'
				}
			});

			if (!response.ok) {
				error = 'Error creating person and relationship';
				return;
			}
			let data = (await response.json()) as {
				person?: components['schemas']['Person'];
				relationships?: components['schemas']['Relationship'][];
			};
			if (onCreation !== undefined) {
				let edges: Array<Edge> = [];
				data.relationships?.map((relationship) =>
					edges.push({
						id: String(relationship.id),
						source: String(relationship.start),
						target: String(relationship.end),
						data: {
							...relationship.properties
						}
					})
				);

				let newNode = {
					id: String(data.person?.Id),
					data: {
						...data.person?.Props
					},
					position: { x: 0, y: 0 },
					type: 'personNode'
				} as Node;
				onCreation([newNode], edges);
				closeModal();
			}
		} else {
			let requestBody = draftPerson as components['schemas']['PersonRegistration'];
			let response = await fetch(`/api/person`, {
				method: 'POST',
				body: JSON.stringify(requestBody),
				headers: {
					'Content-Type': 'application/json'
				}
			});
			if (!response.ok) {
				error = 'Error creating person';
				return;
			}

			if (onOnlyPersonCreation !== undefined) {
				onOnlyPersonCreation(await response.json());
			}
		}

		closeModal();
	}

	onMount(() => {
		if (birth_date) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: birth_date,
					onOpen: function () {
						birth_date.placeholder = '';
					},
					onSelect: function (date) {
						birth_date.value = date.toISOString().split('T')[0];
						draftPerson.born = date.toISOString().split('T')[0];
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
		if (relationship_from_time) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: relationship_from_time,
					onOpen: function () {
						relationship_from_time.placeholder = '';
					},
					onSelect: function (date) {
						relationship_from_time.value = date.toISOString().split('T')[0];
						draftRelationship.from = date.toISOString().split('T')[0];
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
		if (relationship_until) {
			import('pikaday').then(({ default: Pikaday }) => {
				const picker = new Pikaday({
					format: 'YYYY-MM-DD',
					minDate: new Date(1900, 0, 1),
					field: relationship_until,
					onOpen: function () {
						relationship_until.placeholder = '';
					},
					onSelect: function (date) {
						relationship_until.value = date.toISOString().split('T')[0];
						draftRelationship.to = date.toISOString().split('T')[0];
					}
				});
				// Clean up when component unmounts
				return () => picker.destroy();
			});
		}
	});
</script>

<div class="modal modal-open" transition:fade>
	<div
		class="modal-box flex max-h-screen w-full max-w-5xl flex-col items-center justify-center overflow-y-auto"
	>
		<div class="flex w-full max-w-5xl items-center justify-between p-2">
			<h3 class="text-left text-lg font-bold">{create_relationship_and_person()}</h3>
			<div>
				<button class="btn btn-error btn-sm" onclick={onClose}>
					{close()}
				</button>
			</div>
		</div>
		<div class="divider"></div>

		<form onsubmit={onCreate} class="w-full">
			<fieldset
				class="fieldset grid w-full grid-cols-1 items-center gap-y-4 md:grid-cols-2 md:gap-x-6"
			>
				{#if error}
					<div role="alert" class="alert alert-error col-span-full">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 shrink-0 stroke-current"
							fill="none"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
						<span>{error}</span>
					</div>
				{/if}
				{#if relationshipStartID !== undefined}
					<input type="hidden" name="relationshipStartID" value={relationshipStartID} />
					<div class="flex flex-col">
						<label class="label" for="relationship_type">{relation_type()}</label>
						<select
							name="relationship_type"
							class="select select-bordered"
							id="relationship_type"
							bind:value={draftRelationship.type}
						>
							<option value="child">{child()}</option>
							<option value="parent">{parent()}</option>
							<option value="sibling">{sibling()}</option>
							<option value="spouse">{spouse()}</option>
						</select>
					</div>
					<div class="flex flex-col">
						<label class="label" for="relationship_notes"
							>{relation() + ' ' + notes().toLowerCase()}:</label
						>
						<textarea
							name="relationship_notes"
							class="textarea"
							bind:value={draftRelationship.notes}
							placeholder={notes().toLowerCase() + ' ' + optional_field().toLowerCase()}
						></textarea>
					</div>
					<div class="flex flex-col">
						<label class="label" for="relationship_from_time">{from_time()}</label>
						<input
							type="text"
							name="relationship_from_time"
							id="relationship_from_time"
							class="input input-bordered validator pika-single"
							placeholder={optional_field()}
							bind:this={relationship_from_time}
						/>
					</div>
					<div class="flex flex-col">
						<label class="label" for="relationship_until">{until()}</label>
						<input
							type="text"
							name="relationship_until"
							id="relationship_until"
							class="input input-bordered validator pika-single"
							placeholder={optional_field()}
							bind:this={relationship_until}
						/>
					</div>
					<div class="divider margin-t-2 col-span-full"></div>
				{/if}

				<!-- Inputs -->
				<div class="flex flex-col">
					<label class="label" for="first_name">{first_name()}</label>
					<input
						type="text"
						name="first_name"
						class="input input-bordered"
						placeholder={first_name()}
						bind:value={draftPerson.first_name}
					/>
				</div>

				<div class="flex flex-col">
					<label class="label" for="last_name">{last_name()}</label>
					<input
						type="text"
						name="last_name"
						class="input input-bordered"
						placeholder={last_name()}
						bind:value={draftPerson.last_name}
					/>
				</div>

				<div class="flex flex-col">
					<label class="label" for="email">{email()}</label>
					<input
						type="email"
						name="email"
						class="input input-bordered validator"
						placeholder={email() + ' ' + optional_field().toLowerCase()}
						bind:value={draftPerson.email}
					/>
				</div>

				<div class="flex flex-col">
					<label class="label" for="birth_date">{born()}</label>
					<input
						type="text"
						name="birth_date"
						class="input input-bordered validator pika-single"
						placeholder={born()}
						bind:this={birth_date}
					/>
				</div>

				<div class="flex flex-col">
					<label class="label" for="biological_sex">{biological_sex()}</label>
					<select
						name="biological_sex"
						class="select select-bordered"
						id="biological_sex"
						bind:value={draftPerson.biological_sex}
					>
						<option value="male">{male()}</option>
						<option value="female">{female()}</option>
						<option value="intersex">{intersex()}</option>
						<option value="other">{other()}</option>
					</select>
				</div>

				<div class="flex flex-col">
					<label class="label" for="mothers_last_name">{mothers_last_name()}</label>
					<input
						type="text"
						name="mothers_last_name"
						class="input input-bordered"
						placeholder={mothers_last_name()}
						bind:value={draftPerson.mothers_last_name}
					/>
				</div>

				<div class="flex flex-col">
					<label class="label" for="mothers_first_name">{mothers_first_name()}</label>
					<input
						type="text"
						name="mothers_first_name"
						class="input input-bordered"
						placeholder={mothers_first_name()}
						bind:value={draftPerson.mothers_first_name}
					/>
				</div>

				<!-- Submit button spans full width -->
				<div class="col-span-full mt-4 flex justify-center">
					<button type="submit" class="btn btn-neutral mt-4">{create()}</button>
				</div>
			</fieldset>
		</form>
	</div>
</div>
