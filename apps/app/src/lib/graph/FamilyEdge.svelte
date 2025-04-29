<script lang="ts">
	import type { components } from '$lib/api/api.gen.ts';
	import { child, spouse, parent, sibling } from '$lib/paraglide/messages';
	import { getSmoothStepPath, BaseEdge, EdgeLabelRenderer, type EdgeProps } from '@xyflow/svelte';

	let {
		sourceX,
		sourceY,
		source,
		sourcePosition,
		target,
		targetX,
		targetY,
		targetPosition,
		markerEnd,
		style,
		data
	}: EdgeProps = $props();

	let [edgePath, labelX, labelY] = $derived(
		getSmoothStepPath({
			sourceX,
			sourceY,
			sourcePosition,
			targetX,
			targetY,
			targetPosition
		})
	);

	const onEdgeClick = () => {
		window.dispatchEvent(
			new CustomEvent('edge-click', {
				detail: {
					start: source,
					end: target,
					data: data as components['schemas']['FamilyRelationship'] & { type: string }
				}
			})
		);
	};

	let edgeType = (
		data as components['schemas']['FamilyRelationship'] & { type: string }
	).type.toLowerCase();
	let edgeLabel: string = $state(edgeType);
	let edgeColor: string;
	if (edgeType === 'spouse') {
		edgeColor = 'stroke: red;';
		edgeLabel = spouse();
	} else if (edgeType === 'child') {
		edgeColor = 'stroke: blue;';
		edgeLabel = child();
	} else if (edgeType === 'parent') {
		edgeColor = 'stroke: green;';
		edgeLabel = parent();
	} else if (edgeType === 'sibling') {
		edgeColor = 'stroke: brown;';
		edgeLabel = sibling();
	} else {
		edgeColor = 'stroke: gray;';
		edgeLabel = edgeType;
	}
</script>

<BaseEdge path={edgePath} {markerEnd} {style} />
<EdgeLabelRenderer>
	<div
		class="button-edge__label nodrag nopan"
		style:transform="translate(-50%, -50%) translate({labelX}px,{labelY}px)"
	>
		<button class="button-edge__button" onclick={onEdgeClick}>{edgeLabel}</button>
	</div>
</EdgeLabelRenderer>
