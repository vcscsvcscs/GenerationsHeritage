MATCH (p:Person)-[l:Likes]->(r:Recipe)
WHERE id(p) = $personId AND id(r) = $recipeId
DELETE l