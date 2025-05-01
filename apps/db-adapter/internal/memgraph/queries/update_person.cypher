MATCH (n:Person)
WHERE id(n) = $id
SET n += $props
RETURN n AS person