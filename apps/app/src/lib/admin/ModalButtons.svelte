<script lang="ts">
	import {
		add_relationship,
		back,
		biography,
		close,
		create,
		create_person,
		edit,
		managed_profiles,
		relation,
		save
	} from '$lib/paraglide/messages';
	export let createProfile: () => void;
	export let onClose: () => void;
</script>

<div class="flex items-center justify-between p-2">
	<h3 class="text-lg font-bold">{managed_profiles()}</h3>
	<div class="space-x-2">
		<button class="btn btn-success btn-sm" on:click={createProfile}>
			{'+ ' + create_person()}
		</button>
		<button class="btn btn-error btn-sm" on:click={onClose}>
			{close()}
		</button>
	</div>
</div>
