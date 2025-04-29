<script>
	import { deleted_profiles, managed_profiles } from '$lib/paraglide/messages';
</script>

<div class="drawer-side">
	<label for="my-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
	<ul class="menu bg-base-200 text-base-content min-h-full w-80 p-4">
		<!-- Sidebar content here -->
		<li><button class="btn btn-soft btn-primary">{managed_profiles}</button></li>
		<li><button class="btn btn-soft btn-primary">{deleted_profiles()}</button></li>
	</ul>
</div>
