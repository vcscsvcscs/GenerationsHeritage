MATCH (r:Recipe)
WHERE id(r) = $id
SET r += $props
RETURN r AS recipe