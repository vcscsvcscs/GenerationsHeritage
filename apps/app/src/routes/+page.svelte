<script lang="ts">
    import {title, family_tree} from '$lib/paraglide/messages.js';
    import { SvelteFlowProvider,Background, BackgroundVariant, SvelteFlow, Controls, MiniMap } from '@xyflow/svelte';
    import type { Node, Edge, NodeTypes, NodeProps } from '@xyflow/svelte';

    let nodes = $state.raw<Node[]>([]);
    let edges = $state.raw<Edge[]>([]);
</script>

<svelte:head>
	<title>{title({ page: family_tree() })}</title>
</svelte:head>

<div style="height:100vh;">
	<SvelteFlowProvider>
		<SvelteFlow bind:nodes bind:edges class="bg-base-100" fitView onlyRenderVisibleElements>
			<Controls class="bg-base-300 text-base-content" />
			<MiniMap class="bg-base-200" />
		</SvelteFlow>
	</SvelteFlowProvider>
</div>

