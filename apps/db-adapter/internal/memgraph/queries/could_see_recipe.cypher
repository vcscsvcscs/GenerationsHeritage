MATCH (creator:Person)-[:Created]->(r:Recipe)
WHERE id(r) = $recipeId
WITH creator, r
MATCH (viewer:Person)
WHERE id(viewer) = $userId
WITH creator, r, viewer
WHERE id(creator) = id(viewer)
   OR EXISTS {
      MATCH (viewer)-[:Admin]->(creator)
   }
   OR EXISTS {
      MATCH (viewer)-[:Child|Parent|Sibling|Spouse]-(creator)
   }
RETURN r
LIMIT 1
