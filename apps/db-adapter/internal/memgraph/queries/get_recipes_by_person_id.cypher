MATCH (p:Person)
WHERE id(p) = $id
MATCH (viewer:Person)
WHERE id(viewer) = $userId
MATCH (p)-[:Likes|Created]->(r:Recipe)
WITH DISTINCT viewer, p, r
OPTIONAL MATCH (p)-[l:Likes]->(r)
OPTIONAL MATCH (p)-[c:Created]->(r)
OPTIONAL MATCH (creator:Person)-[:Created]->(r)
OPTIONAL MATCH (viewer)-[adm:Admin]->(creator)
WITH r, l, c, viewer, creator, adm
ORDER BY id(r)
RETURN collect({
  recipe: r,
  relationship: l,
  created: c IS NOT NULL,
  can_edit: CASE WHEN id(creator) = id(viewer) OR adm IS NOT NULL THEN true ELSE false END
}) AS entries
