MATCH (n:DeletedPerson {id: $id})
DETACH DELETE n;