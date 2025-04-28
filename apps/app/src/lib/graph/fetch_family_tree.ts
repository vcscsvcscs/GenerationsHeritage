import type { components } from '$lib/api/api.gen';
import type { Layout } from '$lib/graph/model';
import type { Edge, Node } from '@xyflow/svelte';

export function parseFamilyTree(data: components['schemas']['FamilyTree']): Layout {
	
	if (
		data === null ||
		data?.people === null ||
		data?.people === undefined ||
		data?.people.length === 0
	) {
		throw new Error('Family tree is empty');
	}

	const nodes: Node[] = data.people.map((person) => {
		let newNode = { data: { ...person } } as Node;
		if (person.id !== null && person.id !== undefined) {
			newNode.id = person.id.toString();
		}

		return newNode;
	});

	let relationships: Edge[] = [];
	if (data.relationships) {
		relationships = data.relationships.map((relationship) => {
			let newEdge = { data: { ...relationship.properties } } as Edge;
			if (relationship.start !== null && relationship.start !== undefined) {
				newEdge.source = relationship.start.toString();
			}
			if (relationship.end !== null && relationship.end !== undefined) {
				newEdge.target = relationship.end.toString();
			}

			return newEdge;
		});
	}

	return { Nodes: nodes, Edges: relationships };
}
