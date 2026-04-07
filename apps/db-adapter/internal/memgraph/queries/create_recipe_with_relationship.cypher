MATCH (p:Person)
WHERE id(p) = $personId
CREATE (r:Recipe $Recipe)
CREATE (p)-[:Created]->(r)
CREATE (p)-[l:Likes $Relationship]->(r)
RETURN r as recipe, l as relationship