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
			newNode.id = "help"+person.id.toString();
		}
		newNode.data.id = person.id;
		return newNode;
	});

	let relationships: Edge[] = [];
	if (data.relationships) {
		relationships = data.relationships.map((relationship) => {
			const newEdge = { data: { ...relationship.Props } } as Edge;
			newEdge.data!.type = relationship.Type?.toLowerCase();
			if (relationship.StartElementId !== null && relationship.StartElementId !== undefined) {
				newEdge.source = "help"+relationship.StartId!.toString();
			}
			if (relationship.EndElementId !== null && relationship.EndElementId !== undefined) {
				newEdge.target = "help"+relationship.EndId!.toString();
			}

			return newEdge;
		});
	}

	return { Nodes: nodes, Edges: relationships };
}
