MATCH (a)-[r1:Admin]->(b)
WHERE id(a) = $id1 AND id(b) = $id2 AND $id1 != $id2
DELETE r1;