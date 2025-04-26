MATCH (a)-[r1:Admin]->(b)
WHERE id(a) = $id
RETURN collect({id: id(b), first_name: b.first_name, last_name: b.last_name, adminSince: r1.added}) as managed;