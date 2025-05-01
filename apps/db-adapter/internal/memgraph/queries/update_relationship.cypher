MATCH (n)-[r]->(o)
WHERE id(n) = $id1 AND id(o) = $id2
SET r += $relationship
RETURN collect(r) as relationship
