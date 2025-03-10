<script lang="ts">
    import {sign_in,site_intro,title,family_tree} from '$lib/paraglide/messages.js';
    import Google from "./web_neutral_rd_SI.svg";
    import FamilyTree from "./highresolution_icon_no_background_croped.png";
    import ThemeButton from '$lib/theme-select.svelte';
</script>

<svelte:head>
	<title>{title({ page: sign_in() })}</title>
</svelte:head>

<div class="hero bg-base-200 min-h-screen">
    <div class="hero-content flex flex-col items-center justify-center text-center">
      <div class="max-w-md">
        <figure class="px-10 pt-10 top-margin-10">
          <img
            src="{FamilyTree}"
            alt="{family_tree()}"
            class="rounded-xl" />
        </figure>
        <h1 class="text-5xl font-bold">{sign_in()}</h1>
        <p class="py-6">
          {site_intro()}
        </p>
        <a href="/login/google" class="inline-block">
          <img src="{Google}" alt="{sign_in()}" class="w-full h-auto">
        </a>
      </div>
    </div>
</div>

<div class="absolute top-2 right-2">
  <ThemeButton />
</div>