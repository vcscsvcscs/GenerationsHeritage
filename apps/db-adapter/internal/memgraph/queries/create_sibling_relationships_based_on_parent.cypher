MATCH (a:Person)-[:Parent]->(b:Person)-[:Child]->(c:Person)
WHERE id(a) = $childId AND id(b) = $parentId AND $parentId != $childId AND id(c) != id(a)
MERGE (a)-[r1:Sibling]->(c)-[r2:Sibling]->(a)
RETURN collect(r1)+collect(r2) as relationships;