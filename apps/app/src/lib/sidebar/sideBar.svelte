<script>
	import { managed_profiles } from '$lib/paraglide/messages';
</script>

<div class="drawer-side">
	<label for="my-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
	<ul class="menu bg-base-200 text-base-content min-h-full w-80 p-4 gap-4 pt-16">
		<!-- Sidebar content here -->
		<li><button class="btn btn-primary">{managed_profiles()}</button></li>
	</ul>
</div>
