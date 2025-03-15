<script lang="ts">
	import { sign_in, site_intro, title, family_tree, welcome } from '$lib/paraglide/messages.js';
	import FamilyTree from './highresolution_icon_no_background_croped.png';
</script>

<svelte:head>
	<title>{title({ page: sign_in() })}</title>
</svelte:head>

<div class="hero bg-base-200 min-h-screen">
	<div class="hero-content flex flex-col items-center justify-center text-center">
		<div class="max-w-lg">
			<figure class="top-margin-10 px-10 pt-10">
				<img src={FamilyTree} alt={family_tree()} class="rounded-xl" />
			</figure>
			<h1 class="text-3xl font-bold">{welcome()}</h1>
			<p class="py-6">
				{site_intro()}
			</p>
			<a href="/login/google" class="btn rounded-full bg-white text-black border-[#e5e5e5]">
				<!-- Google -->
				<svg
					aria-label="Google logo"
					width="16"
					height="16"
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 512 512"
					><g
						><path d="m0 0H512V512H0" fill="none"></path><path
							fill="#34a853"
							d="M153 292c30 82 118 95 171 60h62v48A192 192 0 0190 341"
						></path><path fill="#4285f4" d="m386 400a140 175 0 0053-179H260v74h102q-7 37-38 57"
						></path><path fill="#fbbc02" d="m90 341a208 200 0 010-171l63 49q-12 37 0 73"></path>
						<path fill="#ea4335" d="m153 219c22-69 116-109 179-50l55-54c-78-75-230-72-297 55"
						></path></g
					></svg
				>
				Google {sign_in()}
			</a>
		</div>
	</div>
</div>
