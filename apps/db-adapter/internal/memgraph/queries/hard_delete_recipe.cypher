MATCH (r:DeletedRecipe)
WHERE id(r) = $id
DETACH DELETE r;