MATCH (a:Person), (b:Person)
WHERE a.id = $id1 AND b.id = $id2
CREATE (a)-[r:Relationship $Relationship]->(b)
RETURN r as relationship