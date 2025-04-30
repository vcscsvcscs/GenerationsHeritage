<script lang="ts">
	import type { components } from '$lib/api/api.gen.ts';
	import { child, spouse, parent, sibling } from '$lib/paraglide/messages';
	import { getBezierPath, BaseEdge, type EdgeProps, Position } from '@xyflow/svelte';

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


	let edgeType = (
		data as components['schemas']['FamilyRelationship'] & { type: string }
	).type.toLowerCase();
	let edgeLabel: string = $state(edgeType);
	let edgeColor: string;
	let srcPos;
	let tgtPos;
	if (edgeType === 'spouse') {
		edgeColor = 'stroke: red;';
		edgeLabel = spouse();
		if (sourceX < targetX) {
			tgtPos = Position.Right;
			srcPos = Position.Left;
		} else {
			tgtPos = Position.Left;
			srcPos = Position.Right;
		}
	} else if (edgeType === 'child') {
		edgeColor = 'stroke: blue;';
		edgeLabel = child();
		if (sourceY < targetY) {
			tgtPos = Position.Bottom;
			srcPos = Position.Top;
		} else {
			tgtPos = Position.Bottom;
			srcPos = Position.Top;
		}
	} else if (edgeType === 'parent') {
		edgeColor = 'stroke: green;';
		edgeLabel = parent();
		if (sourceY < targetY) {
			tgtPos = Position.Bottom;
			srcPos = Position.Top;
		} else {
			tgtPos = Position.Bottom;
			srcPos = Position.Top;
		}
	} else if (edgeType === 'sibling') {
		edgeColor = 'stroke: brown;';
		edgeLabel = sibling();
		if (sourceX < targetX) {
			tgtPos = Position.Right;
			srcPos = Position.Left;
		} else {
			tgtPos = Position.Left;
			srcPos = Position.Right;
		}
	} else {
		edgeColor = 'stroke: gray;';
		edgeLabel = edgeType;
	}

	let  [path, labelX, labelY] = $derived(
		getBezierPath({
			sourceX,
			sourceY,
			sourcePosition: srcPos,
			targetX,
			targetY,
			targetPosition: tgtPos
		})
	);

	edgeColor = edgeColor +'stroke-opacity:unset; stroke-width=20;' +(style ?? '');

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
</script>

<BaseEdge {path} {labelX} {labelY} {markerEnd} style={edgeColor} onclick={onEdgeClick}/>
