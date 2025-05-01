MATCH (a)-[r:Comment]->(b)
WHERE id(a) = $id1 AND id(b) = $id2
DELETE r;