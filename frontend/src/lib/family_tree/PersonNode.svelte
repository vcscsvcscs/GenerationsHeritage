<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	type $$Props = NodeProps;

	export let id: $$Props['id'];
	id;
	export let dragHandle: $$Props['dragHandle'] = undefined;
	dragHandle;
	export let type: $$Props['type'] = undefined;
	type;
	export let selected: $$Props['selected'] = undefined;
	selected;
	export let width: $$Props['width'] = undefined;
	width;
	export let height: $$Props['height'] = undefined;
	height;
	export let dragging: $$Props['dragging'];
	dragging;
	export let targetPosition: $$Props['targetPosition'] = undefined;
	targetPosition;
	export let sourcePosition: $$Props['sourcePosition'] = undefined;
	sourcePosition;
	export let data = {
		ID: '',
		Lastname: 'Nem',
		Firstname: 'Ismert',
		Middlename: '',
		ProfilePicture: 'https://www.flaticon.com/free-icons/user'
	};
</script>

<div class="rounded-badge card card-compact bg-base-300">
	<div class="card-body items-center text-center w-30">
		<div class="avatar">
			<figure class="w-24 mask mask-squircle">
				<img src={data.ProfilePicture} alt="Picture of {data.Lastname} {data.Firstname}" />
			</figure>
		</div>
		<h2 class="card-title text-base-content">
			
			{data.Firstname}
			{data.Middlename}
			{data.Lastname}
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