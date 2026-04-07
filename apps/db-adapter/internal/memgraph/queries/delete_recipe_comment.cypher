MATCH (p:Person)-[c:CommentedOnRecipe]->(r:Recipe)
WHERE id(p) = $personId AND id(r) = $recipeId
DELETE c
