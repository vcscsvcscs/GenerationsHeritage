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
    
    // Determine edge styling and positioning based on relationship type and handles
    if (edgeType === 'spouse') {
        edgeColor = 'stroke: red;';
        edgeLabel = spouse();
        // Use handle-based positioning for spouses
        if (sourceHandleId === 'spouse-right') {
            srcPos = Position.Right;
            tgtPos = Position.Left;
        } else if (sourceHandleId === 'spouse-left') {
            srcPos = Position.Left;
            tgtPos = Position.Right;
        } else {
            // Fallback to position-based logic
            if (sourceX < targetX) {
                srcPos = Position.Right;
                tgtPos = Position.Left;
            } else {
                srcPos = Position.Left;
                tgtPos = Position.Right;
            }
        }
    } else if (edgeType === 'child') {
        edgeColor = 'stroke: blue;';
        edgeLabel = child();
        // Parent-child: from parent's bottom (child handle) to child's top (parent handle)
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
        // Use handle-based positioning for siblings
        if (sourceHandleId === 'spouse-right') {
            srcPos = Position.Right;
            tgtPos = Position.Left;
        } else if (sourceHandleId === 'spouse-left') {
            srcPos = Position.Left;
            tgtPos = Position.Right;
        } else {
            // Fallback to position-based logic
            if (sourceX < targetX) {
                srcPos = Position.Right;
                tgtPos = Position.Left;
            } else {
                srcPos = Position.Left;
                tgtPos = Position.Right;
            }
        }
    } else {
        edgeColor = 'stroke: gray;';
        edgeLabel = edgeType;
        // Keep original positions for unknown types
        srcPos = sourcePosition || Position.Bottom;
        tgtPos = targetPosition || Position.Top;
    }
    
    // Explicit handle overrides (these should take precedence)
    if (sourceHandleId === 'child') {
        srcPos = Position.Bottom;
    }
    if (targetHandleId === 'parent') {
        tgtPos = Position.Top;
    }
    if (sourceHandleId === 'spouse-left') {
        srcPos = Position.Left;
    }
    if (sourceHandleId === 'spouse-right') {
        srcPos = Position.Right;
    }
    if (targetHandleId === 'spouse-left') {
        tgtPos = Position.Left;
    }
    if (targetHandleId === 'spouse-right') {
        tgtPos = Position.Right;
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