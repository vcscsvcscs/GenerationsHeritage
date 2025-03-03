MATCH (n {id: $id1})-[r]-(o {id: $id2})
RETURN r as relationship