MATCH (n:Person)
WHERE id(n) = $id
OPTIONAL MATCH (n)-[p:Parent*..]->(family:Person)
OPTIONAL MATCH (family)-[c:Child*1..4]->(children:Person)
OPTIONAL MATCH (family)-[s:Sibling]->(siblings:Person)
OPTIONAL MATCH (n)-[ds:Sibling]->(direct_siblings:Person)
OPTIONAL MATCH (family)-[fsp:Spouse]->(fspouse:Person)
OPTIONAL MATCH (children)-[csp:Spouse]->(cspouse:Person)
OPTIONAL MATCH (n)-[sp:Spouse]->(spouse:Person)
WITH collections.to_set(collect(n) + collect(family) + collect(children) + collect(direct_siblings) + collect(fspouse) + collect(cspouse) + collect(spouse) + collect(siblings)) as people, 
collections.to_set(collect(c) + collect(p) + collect(s) + collect(ds) + collect(fsp) + collect(csp) + collect(sp)) as relationships
UNWIND people as ppl
RETURN collect({
  id: id(ppl), 
  first_name: ppl.first_name, 
  middle_name: ppl.middle_name,
  last_name: ppl.last_name, 
  born: ppl.born,
  biological_sex: ppl.biological_sex,
  died: ppl.died,
  profile_picture: ppl.profile_picture
}) as people, 
relationships;  