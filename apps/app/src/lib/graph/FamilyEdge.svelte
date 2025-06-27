<script lang="ts">
    import type { components } from '$lib/api/api.gen.ts';
    import { child, spouse, parent, sibling } from '$lib/paraglide/messages';
    import { getSmoothStepPath, BaseEdge, type EdgeProps, Position } from '@xyflow/svelte';
    
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
    let edgeColor: string = $state('stroke: gray;');
    let srcPos: Position = $state(sourcePosition || Position.Bottom);
    let tgtPos: Position = $state(targetPosition || Position.Top);
    
    // Determine edge styling and positioning based on relationship type
    if (edgeType === 'spouse') {
        edgeColor = 'stroke: red;';
        edgeLabel = spouse();
        // Horizontal connection for spouses
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
        // Child relationship: from parent (top) to child (bottom)
        srcPos = Position.Bottom;
        tgtPos = Position.Top;
    } else if (edgeType === 'parent') {
        edgeColor = 'stroke: blue;';
        edgeLabel = parent();
        // Parent relationship: from child (top) to parent (bottom)
        // Note: This might be the inverse of child depending on your data structure
        srcPos = Position.Top;
        tgtPos = Position.Bottom;
    } else if (edgeType === 'sibling') {
        edgeColor = 'stroke: orange;';
        edgeLabel = sibling();
        // Horizontal connection for siblings, but slightly different approach
        const yDiff = Math.abs(sourceY - targetY);
        const xDiff = Math.abs(sourceX - targetX);
        
        if (xDiff > yDiff) {
            // More horizontal separation - use left/right
            if (sourceX < targetX) {
                srcPos = Position.Right;
                tgtPos = Position.Left;
            } else {
                srcPos = Position.Left;
                tgtPos = Position.Right;
            }
        } else {
            // More vertical separation - use top/bottom
            if (sourceY < targetY) {
                srcPos = Position.Bottom;
                tgtPos = Position.Top;
            } else {
                srcPos = Position.Top;
                tgtPos = Position.Bottom;
            }
        }
    } else {
        edgeColor = 'stroke: gray;';
        edgeLabel = edgeType;
        // Keep original positions for unknown types
        srcPos = sourcePosition || Position.Bottom;
        tgtPos = targetPosition || Position.Top;
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