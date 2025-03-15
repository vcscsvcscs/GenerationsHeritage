MATCH (a:Person), (b:Person)
WHERE id(a) = $id1 AND id(b) = $id2
CREATE (a)-[r1:Relationship $Relationship1]->(b)
CREATE (b)-[r2:Relationship $Relationship2]->(a)
RETURN r1 as relationship1, r2 as relationship2 