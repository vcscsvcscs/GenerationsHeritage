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

<div class="rounded-badge card card-compact bg-base-content">
	{#if !isConnecting}
		<Handle class="customHandle" position={Position.Right} type="source" style="z-index: 1;" />
	{/if}
	<!-- <Handle class="customHandle" position={Position.Right} type="source" /> -->

	<Handle class="customHandle" position={Position.Left} type="target" isConnectableStart={false} />
	<div class="card-body w-30 items-center text-center">
		<div class="avatar">
			<figure class="mask mask-squircle w-24">
				<img src={data.profile_picture} alt="Picture of {data.last_name} {data.first_name}" />
			</figure>
		</div>
		<h2 class="card-title text-base-content">
			{data.first_name}
			{data.middle_name ? data.middle_name : ''}
			{data.last_name}
		</h2>
		<h3 class="card-title text-base-content">
			{data.born}
			{data.death ? ' - ' + data.death : ''}
		</h3>
	</div>
</div>
