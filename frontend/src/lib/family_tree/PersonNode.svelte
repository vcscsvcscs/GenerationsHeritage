<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';

	export let id: NodeProps['id'];
	id;
	export let dragHandle: NodeProps['dragHandle'] = undefined;
	dragHandle;
	export let type: NodeProps['type'] = undefined;
	type;
	export let selected: NodeProps['selected'] = undefined;
	selected;
	export let width: NodeProps['width'] = undefined;
	width;
	export let height: NodeProps['height'] = undefined;
	height;
	export let dragging: NodeProps['dragging'];
	dragging;
	export let targetPosition: NodeProps['targetPosition'] = undefined;
	targetPosition;
	export let sourcePosition: NodeProps['sourcePosition'] = undefined;
	sourcePosition;
	export let data = {
		id: '',
		last_name: 'Nem',
		first_name: 'Ismert',
		middle_name: '',
		profile_picture: 'https://cdn-icons-png.flaticon.com/512/3607/3607444.png'
	};
</script>

<div class="rounded-badge card card-compact bg-base-300">
	<div class="card-body items-center text-center w-30">
		<div class="avatar">
			<figure class="w-24 mask mask-squircle">
				<img src={'https://cdn-icons-png.flaticon.com/512/3607/3607444.png'} alt="Picture of {data.last_name} {data.first_name}" />
			</figure>
		</div>
		<h2 class="card-title text-base-content">
			{data.first_name}
			{data.middle_name ? data.middle_name : ''}
			{data.last_name}
		</h2>
	</div>
</div>
<Handle
	type="target"
	position={Position.Top}
	id="spouse"
	style="transform: translate(10px, 50%); left: 0;"
/>
<Handle
	type="source"
	position={Position.Top}
	id="spouse"
	style="transform: translate(0, 50%); left: auto; right: 10px"
/>
<Handle type="target" position={Position.Top} id="parent" style="transform: translate(0, 50%);" />
<Handle
	type="source"
	position={Position.Bottom}
	id="child"
	style="transform: translate(0, -3px);"
/>
<Handle type="target" position={Position.Top} id="child" style="transform: translate(0, 50%);" />
<Handle
	type="source"
	position={Position.Bottom}
	id="parent"
	style="transform: translate(0, -3px);"
/>
