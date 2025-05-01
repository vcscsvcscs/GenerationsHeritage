MATCH (a:Person), (b:Person)
WHERE id(a) = $id1 AND id(b) = $id2 AND $id1 != $id2
MERGE (a)-[r1:Admin]->(b)
ON CREATE SET r1 = {added : timestamp()}
RETURN r1 as relationship;