MATCH (a)-[r1:Admin]->(b)
WHERE id(b) = $id
RETURN collect({id: id(a), first_name: a.first_name, last_name: a.last_name, adminSince: r1.added}) as admins;