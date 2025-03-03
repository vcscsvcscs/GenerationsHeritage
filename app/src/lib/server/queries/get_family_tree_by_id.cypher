MATCH (n:Person { id: $id })-[p:Parent*1..]->(family:Person)
OPTIONAL MATCH (family)-[c:Child]->(children:Person)
WITH family, p, children, c, n
OPTIONAL MATCH (children)<-[p2:Parent]-(OtherParents:Person)
WITH family, p, children, c, OtherParents, p2, n
OPTIONAL MATCH (family)-[s:Spouse]-(spouse:Person)
RETURN family, p, children, c, OtherParents, p2, spouse, s, n;
