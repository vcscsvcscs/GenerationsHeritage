MATCH (a:Person), (b:Person)
WHERE id(a) = $childId AND id(b) = $parentId AND $parentId != $childId
MERGE (a)-[r1:Parent]->(b)-[r2:Child]->(a)
ON CREATE SET r1 = $childRelationship
ON CREATE SET r2 = $parentRelationship
RETURN collect(r1)+collect(r2) as relationships;