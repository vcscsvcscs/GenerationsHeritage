<script>
	import { Handle, Position } from '@xyflow/svelte';

	export let person = {
		ID : "",
		Lastname : "Nem",
		Firstname : "Ismert",
		Middlename : "",
		ProfilePicture : "https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg",
		
	};
</script>

<div class="rounded-badge card card-compact bg-primary">
	<div class="card-body items-center text-center w-30">
		<div class="avatar">
			<figure class="w-24 mask mask-squircle">
				<img src="{person.ProfilePicture}" alt="Picture of {person.Lastname} {person.Firstname}" />
			</figure>
		</div>
		<h2 class="card-title text-primary-content">{person.Lastname} {person.Firstname} {person.Middlename}</h2>
	</div>
</div>
<Handle
	type="target"
	position={Position.Top}
	id="sibling"
	style="transform: translate(10px, 50%); left: 0;"
/>
<Handle
	type="source"
	position={Position.Top}
	id="sibling"
	style="transform: translate(0, 50%); left: auto; right: 10px"
/>
<Handle type="target" position={Position.Top} id="parent" style="transform: translate(0, 50%);" />
<Handle
	type="source"
	position={Position.Right}
	id="right"
	style="transform: translate(-3px, 50%);"
/>
<Handle type="target" position={Position.Left} id="left" style="transform: translate(3px, 50%);" />
<Handle
	type="source"
	position={Position.Bottom}
	id="child"
	style="transform: translate(0, -3px);"
/>
