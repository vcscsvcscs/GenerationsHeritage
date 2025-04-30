<script lang="ts">
	import { managed_profiles } from '$lib/paraglide/messages';
	let clicked = $state(false);
	let {
		open_admin_panel = () => {
			console.log('admin panel opened');
		}
	}: { open_admin_panel: () => void } = $props();
</script>

<div class="dropdown">
	<button
		tabindex="0"
		class={'btn btn-circle swap swap-rotate' + (clicked ? ' swap-active' : '')}
		onclick={() => (clicked = !clicked)}
	>
		<input type="checkbox" />
		<!-- hamburger icon -->
		<svg
			class="swap-off fill-current"
			xmlns="http://www.w3.org/2000/svg"
			width="32"
			height="32"
			viewBox="0 0 512 512"
		>
			<path d="M64,384H448V341.33H64Zm0-106.67H448V234.67H64ZM64,128v42.67H448V128Z" />
		</svg>

		<!-- close icon -->
		<svg
			class="swap-on fill-current"
			xmlns="http://www.w3.org/2000/svg"
			width="32"
			height="32"
			viewBox="0 0 512 512"
		>
			<polygon
				points="400 145.49 366.51 112 256 222.51 145.49 112 112 145.49 222.51 256 112 366.51 145.49 400 256 289.49 366.51 400 400 366.51 289.49 256 400 145.49"
			/>
		</svg>
	</button>
	<ul class="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow">
		<li>
			<button
				tabindex="0"
				class="btn btn-primary"
				aria-label="close sidebar"
				onclick={open_admin_panel}
			>
				{managed_profiles()}
			</button>
		</li>
	</ul>
</div>
