<script lang="ts">
	export let showPanel = true;
	export let data = {
		id: '',
		last_name: 'Nem',
		first_name: 'Ismert',
		middle_name: '',
		profile_picture: 'https://cdn-icons-png.flaticon.com/512/3607/3607444.png'
	};
	let dialog: HTMLDialogElement; // HTMLDialogElement

	$: if (dialog && showPanel) dialog.showModal();
</script>

<dialog bind:this={dialog} class="modal bg-base-300">
	<div class="modal-box w-11/12 max-w-5xl">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
		</form>

		<div class="hero min-h-screen bg-base-200">
			<div class="hero-content flex-col lg:flex-row">
				<div class="avatar max-w-sm rounded-lg shadow-2xl">
					<figure class="w-24 mask mask-squircle">
						<img src={data.profile_picture} alt="Picture of {data.last_name} {data.first_name}" />
					</figure>
				</div>
				<div>
					<h1 class="text-5xl font-bold">
						{data.first_name}
						{data.middle_name}
						{data.last_name}
					</h1>
					<p class="py-6">
						Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda excepturi
						exercitationem quasi. In deleniti eaque aut repudiandae et a id nisi.
					</p>
					<button class="btn btn-primary">Get Started</button>
				</div>
			</div>
		</div>
		<form method="dialog" class="modal-backdrop">
			<button>close</button>
		</form>
	</div>
</dialog>
