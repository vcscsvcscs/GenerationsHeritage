MATCH (a:Person)-[r:Comment]->(b:Person)
WHERE id(b) = $id
RETURN collect(r) as comments, collect({
  id: id(a), 
  first_name: a.first_name, 
  middle_name: a.middle_name,
  last_name: a.last_name, 
  born: a.born,
  died: a.died,
  profile_picture: a.profile_picture
}) as people;