MATCH (viewer:Person)
WHERE id(viewer) = $userId
OPTIONAL MATCH (viewer)-[:Admin]->(managed:Person)
WITH viewer, collect(managed) AS managed
OPTIONAL MATCH (viewer)-[:Parent|Child|Sibling|Spouse*1..%d]-(relative:Person)
WITH viewer, managed, collect(DISTINCT relative) AS relatives
UNWIND [viewer] + managed + relatives AS member
MATCH (member)-[:Likes|Created]->(r:Recipe)
WHERE id(r) = $recipeId
RETURN r
LIMIT 1
