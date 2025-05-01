import type { Connection } from '@xyflow/svelte';
import type { EdgeBase } from '@xyflow/system';

export function isValidConnection(edge: EdgeBase | Connection) {
	if (Number(edge.source) !== Number(edge.target)) {
		return true;
	}

	return false;
}
