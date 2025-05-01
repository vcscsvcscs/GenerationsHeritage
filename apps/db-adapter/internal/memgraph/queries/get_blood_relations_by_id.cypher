MATCH (n:Person)
WHERE id(n) = $id
OPTIONAL MATCH (n)-[p:Parent*..]->(ancestors:Person)
optional MATCH (n)-[cc:Child*..]->(descendants:Person)
OPTIONAL MATCH (ancestors)-[c:Child*1..4]->(children:Person)
OPTIONAL MATCH (ancestors)-[s:Sibling]->(siblings:Person)
OPTIONAL MATCH (n)-[ds:Sibling]->(direct_siblings:Person)
WITH collections.to_set(
  collect(n)+
  collect(descendants)+
  collect(ancestors)+
  collect(children)+
  collect(direct_siblings)+
  collect(siblings)
  ) as people, 
collections.to_set(
  collect(c) +
  collect(cc) + 
  collect(p) + 
  collect(s) + 
  collect(ds)
  ) as relationships
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