MATCH (n:Person)
WHERE id(n)=$id
SET n:DeletedPerson
REMOVE n:Person
RETURN labels(n) AS labels, n AS person