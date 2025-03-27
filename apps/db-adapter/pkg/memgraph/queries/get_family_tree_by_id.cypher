MATCH (n:Person)-[p:Parent*1..]->(family:Person)
WHERE id(n) = $id
OPTIONAL MATCH (family)-[c:Child*..4]->(children:Person)
WITH family, p, children, c, n
OPTIONAL MATCH (children)<-[p2:Parent]-(OtherParents:Person)
WITH family, p, children, c, OtherParents, p2, n
OPTIONAL MATCH (family)-[s:Spouse]-(spouse:Person)
RETURN family, p, children, c, OtherParents, p2, spouse, s, n;
