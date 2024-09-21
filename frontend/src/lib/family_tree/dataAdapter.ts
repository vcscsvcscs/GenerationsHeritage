import { type Node, type Edge } from '@xyflow/svelte';

function AddToNodesData(data: any, i: number, pushToNodesCallback: (arg: Node) => void) {
	if (data[0].Values[i] != null) {
		if (Object.prototype.toString.call(data[0].Values[i]) === '[object Array]') {
			data[0].Values[i].forEach((person: { ElementId: string; Props: {} }) => {
				pushToNodesCallback({
					id: person.ElementId,
					type: 'custom',
					data: person.Props,
					position: { x: 0, y: 0 }
				});
			});
		} else {
			console.log(data[0].Values[i]);
			pushToNodesCallback({
				id: data[0].Values[i].ElementId,
				type: 'custom',
				data: data[0].Values[i].Props,
				position: { x: 0, y: 0 }
			});
		}
	}
}

function AddToEdgesData(data: any, i: number, pushToEdgesCallback: (arg: Edge) => void) {
	if (data[0].Values[i] != null) {
		if (Object.prototype.toString.call(data[0].Values[i]) === '[object Array]') {
			data[0].Values[i].forEach(
				(edge: {
					ElementId: string;
					StartElementId: string;
					EndElementId: string;
					Type: string;
					Props: {};
				}) => {
					pushToEdgesCallback({
						id: edge.ElementId,
						source: edge.StartElementId,
						sourceHandle:
							edge.Type === 'Parent' ? 'parent' : edge.Type === 'Child' ? 'child' : 'spouse',
						target: edge.EndElementId,
						targetHandle:
							edge.Type === 'Parent' ? 'child' : edge.Type === 'Child' ? 'parent' : 'spouse',
						label: edge.Type,
						data: edge.Props,
						zIndex: edge.Type === 'Spouse' ? 1 : 0
					});
				}
			);
		} else {
			console.log(data[0].Values[i]);
			pushToEdgesCallback({
				id: data[0].Values[i].ElementId,
				source: data[0].Values[i].StartElementId,
				sourceHandle:
					data[0].Values[i].Type === 'Parent'
						? 'parent'
						: data[0].Values[i].Type === 'Child'
							? 'child'
							: 'spouse',
				target: data[0].Values[i].EndElementId,
				targetHandle:
					data[0].Values[i].Type === 'Parent'
						? 'child'
						: data[0].Values[i].Type === 'Child'
							? 'parent'
							: 'spouse',
				label: data[0].Values[i].Type,
				data: data[0].Values[i].Props,
				zIndex: data[0].Values[i].Type === 'Spouse' ? 1 : 0
			});
		}
	}
}

export { AddToNodesData, AddToEdgesData };
