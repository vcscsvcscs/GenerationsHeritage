MATCH (p:Person)-[:Created]->(r)
WHERE id(r) = $recipeId AND id(p) = $userId
RETURN r
UNION
MATCH (admin:Person)-[:Admin]->(p:Person)-[:Created]->(r)
WHERE id(r) = $recipeId AND id(admin) = $userId
RETURN r
