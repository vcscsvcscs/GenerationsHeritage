<!-- <svelte:options immutable /> -->

<script lang="ts">
	import { death } from './../paraglide/messages.js';
	import { Handle, Position, useConnection, type NodeProps } from '@xyflow/svelte';
	import type { components } from '$lib/api/api.gen';
	type $$Props = NodeProps;

	export let id: NodeProps['id']; 
	export let data: NodeProps['data'] & components['schemas']['PersonProperties'];
    console.log('data', data);
	const connection = useConnection(); 
	let isConnecting = false;
	let isTarget = false;

	$: isConnecting = connection.current.fromHandle !== null;
	$: isTarget = connection.current.toHandle?.id !== id;
</script>

<div class="card card-compact bg-primary-content text-primary rounded-full w-40 h-40 flex flex-col items-center justify-center shadow-lg">
	{#if !isConnecting}
		<Handle class="customHandle" position={Position.Right} type="source" style="z-index: 1;" />
	{/if}
	<Handle class="customHandle" position={Position.Left} type="target" isConnectableStart={false} />

	<div class="avatar mb-2">
        {#if isConnecting && isTarget}
            <Handle position={Position.Left} type="target" isConnectableStart={false} style="z-index: 1;" />
        {/if}
		<div class="w-24 rounded-full ring ring-accent ring-offset-accent ring-offset-1 border-0 bg-accent">
            {#if isConnecting && isTarget}
                <Handle position={Position.Left} type="target" isConnectableStart={false} style="z-index: 1;" />
            {/if}
			<img src={data.profile_picture||'https://cdn-icons-png.flaticon.com/512/10628/10628885.png'} alt="Picture of {data.last_name} {data.first_name}" />
		</div>
	</div>

	<div class="text-center px-2">
		<h2 class="font-semibold text-sm leading-tight">
			{data.first_name} {data.middle_name ? data.middle_name : ''} {data.last_name}
		</h2>
		<h3 class="text-xs opacity-70">
			{data.born}{data.death ? ' - ' + data.death : ''}
		</h3>
	</div>
</div>

<style>
    :global(div.customHandle) {
    width: 100%;
    height: 100%;
    background: blue;
    position: absolute;
    top: 0;
    left: 0;
    border-radius: 0;
    transform: none;
    border: none;
    opacity: 0;
  }
</style>