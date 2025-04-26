import type { components } from '$lib/api/api.gen';
import type { Node } from '@xyflow/svelte';

export function handleNodeClick(set_panel_options:(person:components['schemas']['PersonProperties']&{id:number})=>void): ({
    detail: { event, node }
}: {
    detail: { event: MouseEvent; node: Node };
}) => void {
    return ({ detail: { event, node } }) => {
        event.preventDefault();
        node.data.id = Number(node.id);
        set_panel_options(node.data as components['schemas']['PersonProperties']&{id:number});
    }
}