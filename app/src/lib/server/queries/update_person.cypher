MATCH (n:Person $props)
SET n += $props
RETURN n AS person