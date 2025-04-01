MATCH (n:Person)-[p:Parent*1..]->(family:Person)
WHERE id(n) = $id
OPTIONAL MATCH (family)-[c:Child*1..4]->(children:Person)
RETURN collect(family) + collect(children) + collect(n), collect(c) + collect(p)
