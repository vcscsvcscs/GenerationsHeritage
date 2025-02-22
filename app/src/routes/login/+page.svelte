<script lang="ts">
    import {sign_in,title} from '$lib/paraglide/messages.js';
    import Google from "./web_neutral_rd_SI.svg";
</script>

<svelte:head>
	<title>{title({ page: sign_in() })}</title>
</svelte:head>

<div class="hero bg-base-200 min-h-screen">
    <div class="hero-content flex flex-col items-center justify-center text-center">
      <div class="max-w-md">
        <h1 class="text-5xl font-bold">{sign_in()}</h1>
        <p class="py-6">
          Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda excepturi exercitationem
          quasi. In deleniti eaque aut repudiandae et a id nisi.
        </p>
        <a href="/login/google" class="inline-block">
          <img src="{Google}" alt="{sign_in()}" class="w-full h-auto">
        </a>
      </div>
    </div>
</div>