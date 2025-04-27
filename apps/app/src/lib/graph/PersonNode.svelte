<!-- <svelte:options immutable /> -->

<script lang="ts">
	import { Handle, Position, useConnection, type NodeProps } from '@xyflow/svelte';
	import type { components } from '$lib/api/api.gen';
	import { isValidConnection } from './connection.js';
	type $$Props = NodeProps;

	export let id: NodeProps['id'];
	export let data: NodeProps['data'] & components['schemas']['PersonProperties'];
	const connection = useConnection();
	let isConnecting = false;
	let isTarget = false;

	$: isConnecting = connection.current.fromHandle !== null;
	$: isTarget = connection.current.toHandle?.id !== id;
</script>

<div
	class="card card-compact bg-primary-content text-primary flex h-40 w-40 flex-col items-center justify-center rounded-full shadow-lg"
>
	{#if !isConnecting}
		<Handle class="customHandle" position={Position.Right} type="source" style="z-index: 1;" />
	{/if}
	<Handle class="customHandle" position={Position.Left} type="target" isConnectableStart={false} />

	<div class="avatar mb-2">
		{#if isConnecting && isTarget}
			<Handle
				isValidConnection={isValidConnection}
				position={Position.Left}
				type="target"
				isConnectableStart={false}
				style="z-index: 1;"
			/>
		{/if}
		<div
			class="ring-accent ring-offset-accent bg-accent w-24 rounded-full border-0 ring ring-offset-1"
		>
			{#if isConnecting && isTarget}
				<Handle
					isValidConnection={isValidConnection}
					position={Position.Left}
					type="target"
					isConnectableStart={false}
					style="z-index: 1;"
				/>
			{/if}
			<img
				src={data.profile_picture || 'https://cdn-icons-png.flaticon.com/512/10628/10628885.png'}
				alt="Picture of {data.last_name} {data.first_name}"
			/>
		</div>
	</div>

	<div class="px-2 text-center">
		<h2 class="text-sm leading-tight font-semibold">
			{data.first_name}
			{data.middle_name ? data.middle_name : ''}
			{data.last_name}
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
