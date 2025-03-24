MATCH (a:Person), (b:Person)
WHERE id(a) = $id1 AND id(b) = $id2
CREATE (a)-[r:Relationship $Relationship]-(b)
RETURN r as relationship