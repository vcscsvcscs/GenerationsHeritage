MATCH (p:Person), (r:Recipe)
WHERE id(p) = $personId AND id(r) = $recipeId
MERGE (p)-[c:CommentedOnRecipe]->(r)
SET c += $Comment
RETURN c as comment, p as commenter
