import type { Connection } from '@xyflow/svelte';
import type { EdgeBase } from '@xyflow/system';

export function isValidConnection(edge: EdgeBase | Connection) {
	if (edge.source !== edge.target) {
		return true;
	}

	return false;
}
