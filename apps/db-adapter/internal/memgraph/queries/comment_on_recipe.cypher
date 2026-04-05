MATCH (p:Person), (r:Recipe)
WHERE id(p) = $personId AND id(r) = $recipeId
CREATE (p)-[c:CommentedOnRecipe $Comment]->(r)
RETURN c as comment, p as commenter
