MATCH (a:Person), (b:Person)
WHERE id(a) = $id1 AND id(b) = $id2 AND $id1 != $id2
MERGE (a)-[r1:Sibling]->(b)-[r2:Sibling]->(a)
ON CREATE SET r1 = $Relationship1
ON CREATE SET r2 = $Relationship2
RETURN collect(r1)+collect(r2) as relationships;