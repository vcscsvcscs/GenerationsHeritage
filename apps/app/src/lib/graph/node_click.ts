import type { components } from '$lib/api/api.gen';
import type { NodeEventWithPointer } from '@xyflow/svelte';

export function handleNodeClick(
	set_panel_options: (person: components['schemas']['PersonProperties'] & { id: number }) => void
): NodeEventWithPointer<MouseEvent | TouchEvent> {
	return ({ event, node }) => {
		event.preventDefault();
		set_panel_options(node.data as components['schemas']['PersonProperties'] & { id: number });
	};
}
