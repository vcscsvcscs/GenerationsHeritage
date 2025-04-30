<script lang="ts">
	import {
		hard_delete,
		managed_profiles,
		delete_profile,
		edit,
		from_time,
		admin,
		create_relationship_and_person,

		add_relationship

	} from '$lib/paraglide/messages';
	import ModalButtons from './ModalButtons.svelte';
	import type { components, operations } from '$lib/api/api.gen';

	let {
		closeModal,
		editProfile = () => {},
		onChange = () => {},
		addRelationship = () => {},
		createProfile = ()=> {},
		createRelationshipAndProfile = () => {},
	} = $props<{
		closeModal: () => void;
		onChange?: () => void;
		addRelationship?: (id: number) => void;
		createRelationshipAndProfile?: (id: number) => void;
		editProfile?: (id: number) => void;
		createProfile?: () => void;
	}>();

	let managed_profiles_list: components['schemas']['Admin'][] = $state([]);
	function fetchManagedProfiles(){
	fetch(`/api/managed_profiles`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then((response) => {
			if (!response.ok) {
				console.log('Cannot get managed profiles, status: ' + response.status);
				return;
			}

			return response.json();
		})
		.then((data) => {
			if (data) {
				managed_profiles_list = [...(data as components['schemas']['Admin'][])];
			}
		})
		.catch((error) => {
			console.error('Error fetching managed profiles:', error);
		});
	}

	fetchManagedProfiles();

	async function deleteProfile(id: number) {
		fetch('/api/person/' + id, {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json'
			}
		})
			.then((response) => {
				if (response.ok) {
					onChange();
					managed_profiles_list = managed_profiles_list.filter((profile) => profile.id !== id);
					return;
				} else {
					alert('Error deleting person');
				}
			})
			.catch((error) => {
				console.info('Error:', error);
			});
	}

	async function hardDeleteProfile(id: number) {
		fetch('/api/person/' + id + '/hard-delete', {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json'
			}
		})
			.then((response) => {
				if (response.ok) {
					onChange();
					managed_profiles_list = managed_profiles_list.filter((profile) => profile.id !== id);
					return;
				} else {
					alert('Error deleting person');
				}
			})
			.catch((error) => {
				console.error('Error:', error);
			});
	}
</script>

<div class="modal modal-open z-8">
	<div class="modal-box w-full max-w-xl gap-4">
		<div class="bg-base-100 sticky top-0 z-5">
			<ModalButtons onClose={closeModal} {createProfile} />
			<div class="divider"></div>
		</div>
		<ul class="list bg-base-100 rounded-box shadow-md">
			<li class="p-4 pb-2 text-xs tracking-wide opacity-60">{managed_profiles()}</li>
			{#each managed_profiles_list as profile}
				<li class="list-row">
					<div>
						<div>{profile.first_name + ' ' + profile.last_name}</div>
						<div class="text-xs font-semibold uppercase opacity-60">{profile.id}</div>
					</div>
					<div>
						<div class="text-xs font-semibold uppercase opacity-60">
							{admin() + ' ' + from_time().toLowerCase() + ': ' + profile.adminSince}
						</div>
						<div class="text-xs font-semibold uppercase opacity-60">{profile.label![0]}</div>
					</div>
					<button
					class="btn btn-success btn-soft"
					onclick={() => {
						addRelationship(profile.id!);
					}}>
					{add_relationship()}
				</button>
					<button
					class="btn btn-success btn-soft"
					onclick={() => {
						createRelationshipAndProfile(profile.id!);
					}}>
					{create_relationship_and_person()}
				</button>
					<button
						class="btn btn-secondary"
						onclick={() => {
							editProfile(profile.id!);
						}}>
						{edit()}
					</button>
					{#if profile.label?.includes('DeletedPerson')}
						<button
							class="btn btn-error"
							onclick={() => {
								hardDeleteProfile(profile.id!);
							}}
						>
							{hard_delete()}
						</button>
					{:else}
						<button
							class="btn btn-error"
							onclick={() => {
								deleteProfile(profile.id!);
							}}
						>
							{delete_profile()}
						</button>
					{/if}
				</li>
			{/each}
		</ul>
	</div>
</div>
