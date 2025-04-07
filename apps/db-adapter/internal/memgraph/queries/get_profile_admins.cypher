MATCH (a)-[r1:Admin]->(b)
WHERE id(b) = $id
RETURN collect(r1) as adminRelationship, collect({id: id(a), first_name: a.first_name, last_name: a.last_name}) as admins;