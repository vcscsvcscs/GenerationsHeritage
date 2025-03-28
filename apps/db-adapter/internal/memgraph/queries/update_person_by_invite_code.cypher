MATCH (n:Person)
WHERE n.invite_code = $invite_code
SET n += $props
RETURN n AS person