MATCH (r:Recipe)
WHERE id(r) = $id
SET r:DeletedRecipe
REMOVE r:Recipe
RETURN labels(r) AS labels, r AS recipe