MATCH (n:Person {id: $id})
SET n:DeletedPerson
REMOVE n:Person
RETURN labels(n) AS labels, n AS person