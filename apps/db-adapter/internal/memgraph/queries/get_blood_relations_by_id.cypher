MATCH (n:Person)-[p:Parent*1..]->(family:Person)
WHERE id(n) = $id
OPTIONAL MATCH (family)-[c:Child*1..]->(children:Person)
RETURN family, p, children, c, n
