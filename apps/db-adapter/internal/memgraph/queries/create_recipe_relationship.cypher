MATCH (p:Person), (r:Recipe)
WHERE id(p) = $personId AND id(r) = $recipeId
MERGE (p)-[l:Likes]->(r)
SET l += $Relationship
RETURN l as relationship