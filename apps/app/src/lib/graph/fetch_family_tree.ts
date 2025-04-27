import type {Layout} from '$lib/graph/model';


export async function fetchFamilyTree(with_out_spouse: boolean): Promise<Layout> {
    const url = with_out_spouse
        ? '/api/family_tree?with_out_spouse=true'
        : '/api/family_tree?with_out_spouse=false';
    const response = await fetch(url, {
        method: 'GET'
    })

    if (response.status !== 200){
        throw new Error(await response.text());
    }

    let layout: Layout = {
        Nodes: [],
        Edges: []
    }

    return layout;
}