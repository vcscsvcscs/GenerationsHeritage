<!-- <svelte:options immutable /> -->

<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	import type { components } from '$lib/api/api.gen';
	import { isValidConnection } from './connection.js';
	type $$Props = NodeProps;

	export let data: NodeProps['data'] & components['schemas']['PersonProperties'];

	let nodeColor = ' bg-neutral text-neutral-content';
	switch (data.biological_sex) {
		case 'female':
			nodeColor = ' bg-secondary text-secondary-content';
			break;
		case 'male':
			nodeColor = ' bg-primary text-primary-content';
			break;
		case 'intersex':
			nodeColor = ' bg-accent text-accent-content';
			break;
	}
</script>

<div
	class={'card card-compact flex h-40 w-40 flex-col items-center justify-center rounded-full shadow-lg' +
		nodeColor}
>
	<Handle
		class="customHandle"
		id="child"
		{isValidConnection}
		isConnectable={true}
		position={Position.Bottom}
		type="source"
		style="z-index: 1;"
	/>

	<Handle
		class="customHandle"
		{isValidConnection}
		position={Position.Left}
		isConnectable={true}
		type="target"
		isConnectableStart={false}
	/>

	<Handle
		class="customHandle"
		{isValidConnection}
		position={Position.Right}
		isConnectable={true}
		type="target"
		isConnectableStart={false}
	/>
	<Handle
		class="customHandle"
		{isValidConnection}
		position={Position.Left}
		isConnectable={true}
		type="source"
		isConnectableStart={true}
	/>

	<Handle
		class="customHandle"
		{isValidConnection}
		position={Position.Right}
		isConnectable={true}
		type="source"
		isConnectableStart={true}
	/>
	<Handle
		class="customHandle"
		id="parent"
		{isValidConnection}
		position={Position.Top}
		isConnectable={true}
		type="target"
		isConnectableStart={false}
	/>

	<div class="avatar mb-2" style="z-index: 2; cursor: pointer;">
		<div class="bg-accent w-24 rounded-full border-0 ring-offset-1">
			<img
				src={data.profile_picture || 'https://cdn-icons-png.flaticon.com/512/10628/10628885.png'}
				alt="Picture of {data.last_name} {data.first_name}"
			/>
		</div>
	</div>

	<div class="px-2 text-center" style="z-index: 2; cursor: pointer;">
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
