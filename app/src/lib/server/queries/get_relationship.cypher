MATCH (n)-[r]-(o)
WHERE id(n) = $id1 AND id(o) = $id2
RETURN r as relationship