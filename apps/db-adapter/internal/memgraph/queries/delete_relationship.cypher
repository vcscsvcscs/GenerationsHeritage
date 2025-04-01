MATCH (a:Person)-[r]-(b:Person)
WHERE id(a) = $id1 AND id(b) = $id2 AND type(r) != Admin
DELETE r;