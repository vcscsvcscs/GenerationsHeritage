MATCH (r:Recipe)
WHERE id(r) = $recipeId
MATCH (viewer:Person)
WHERE id(viewer) = $userId
OPTIONAL MATCH (creator:Person)-[:Created]->(r)
OPTIONAL MATCH (viewer)-[adm:Admin]->(creator)
RETURN r AS recipe,
  CASE WHEN id(creator) = id(viewer) OR adm IS NOT NULL THEN true ELSE false END AS can_edit
LIMIT 1
