MATCH (n:DeletedPerson)
WHERE id(n) = $id
DETACH DELETE n;