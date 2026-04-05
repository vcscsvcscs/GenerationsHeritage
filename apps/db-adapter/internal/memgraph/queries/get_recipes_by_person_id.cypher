MATCH (p:Person)-[l:Likes]->(r:Recipe)
WHERE id(p) = $id
RETURN collect(r) as recipes, collect(l) as recipeRelations