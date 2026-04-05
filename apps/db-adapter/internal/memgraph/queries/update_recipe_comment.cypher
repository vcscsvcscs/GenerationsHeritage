MATCH (p:Person)-[c:CommentedOnRecipe]->(r:Recipe)
WHERE id(p) = $personId AND id(r) = $recipeId
SET c.message = $message,
    c.edited = timestamp()
RETURN c as comment
