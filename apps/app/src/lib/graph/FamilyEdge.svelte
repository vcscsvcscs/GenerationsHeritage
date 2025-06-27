<script lang="ts">
    import type { components } from '$lib/api/api.gen.ts';
    import { child, spouse, parent, sibling } from '$lib/paraglide/messages';
    import { getSmoothStepPath, BaseEdge, type EdgeProps, Position } from '@xyflow/svelte';
    
    let {
        sourceX,
        sourceY,
        source,
        sourcePosition,
        sourceHandleId,
        target,
        targetX,
        targetY,
        targetPosition,
        targetHandleId,
        markerEnd,
        style,
        data
    }: EdgeProps = $props();
    
    let edgeType = (
        data as components['schemas']['FamilyRelationship'] & { type: string }
    ).type.toLowerCase();
    
    let edgeLabel: string = $state(edgeType);
    let edgeColor: string = $state('stroke: gray;');
    let srcPos: Position = $state(sourcePosition || Position.Bottom);
    let tgtPos: Position = $state(targetPosition || Position.Top);
    
    // Determine edge styling and positioning based on relationship type
    if (edgeType === 'spouse') {
        edgeColor = 'stroke: red;';
        edgeLabel = spouse();
        // For spouse connections, use horizontal positioning
        if (sourceX < targetX) {
            srcPos = Position.Right;
            tgtPos = Position.Left;
        } else {
            srcPos = Position.Left;
            tgtPos = Position.Right;
        }
    } else if (edgeType === 'child') {
        edgeColor = 'stroke: blue;';
        edgeLabel = child();
        // Use the handles set by the layout: child handle (bottom) to parent handle (top)
        srcPos = Position.Bottom;
        tgtPos = Position.Top;
    } else if (edgeType === 'parent') {
        edgeColor = 'stroke: blue;';
        edgeLabel = parent();
        // Parent relationship: from child (top) to parent (bottom)
        srcPos = Position.Top;
        tgtPos = Position.Bottom;
    } else if (edgeType === 'sibling') {
        edgeColor = 'stroke: orange;';
        edgeLabel = sibling();
        // For siblings, use horizontal positioning
        if (sourceX < targetX) {
            srcPos = Position.Right;
            tgtPos = Position.Left;
        } else {
            srcPos = Position.Left;
            tgtPos = Position.Right;
        }
    } else {
        edgeColor = 'stroke: gray;';
        edgeLabel = edgeType;
        // Keep original positions for unknown types
        srcPos = sourcePosition || Position.Bottom;
        tgtPos = targetPosition || Position.Top;
    }
    
    // Override with handle-specific positioning if handles are specified
    if (sourceHandleId === 'child') {
        srcPos = Position.Bottom;
    }
    if (targetHandleId === 'parent') {
        tgtPos = Position.Top;
    }
    
    let [path, labelX, labelY] = $derived(
        getSmoothStepPath({
            sourceX,
            sourceY,
            sourcePosition: srcPos,
            targetX,
            targetY,
            targetPosition: tgtPos,
            borderRadius: 20, // Add some rounding to make paths smoother
            offset: 20 // Add offset to avoid overlapping with nodes
        })
    );
    
    // Fix the style string formatting
    const finalStyle = `${edgeColor} stroke-width: 3; stroke-opacity: 0.8; ${style ?? ''}`;
    
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

<BaseEdge {path} {labelX} {labelY} {markerEnd} style={finalStyle} onclick={onEdgeClick} />