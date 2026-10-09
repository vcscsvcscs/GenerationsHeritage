MATCH (r)
WHERE id(r) = $id
RETURN labels(r) AS labels
