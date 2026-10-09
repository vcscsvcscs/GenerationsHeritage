MATCH (n:Person)
WHERE id(n) = $id
OPTIONAL MATCH (n)-[:Parent|Child|Sibling|Spouse*1..%d]-(relative:Person)
WITH n, collect(DISTINCT relative) + [n] AS family
UNWIND family AS member
MATCH (member)-[:Likes|Created]->(r:Recipe)
OPTIONAL MATCH (member)-[l:Likes]->(r)
OPTIONAL MATCH (creator:Person)-[:Created]->(r)
OPTIONAL MATCH (n)-[adm:Admin]->(creator)
WITH n, r, member, l,
  CASE WHEN id(member) = id(n) THEN 0 WHEN id(member) = id(creator) THEN 1 ELSE 2 END AS priority,
  CASE WHEN id(creator) = id(n) OR adm IS NOT NULL THEN true ELSE false END AS can_edit
ORDER BY id(r), priority, id(member)
WITH r, can_edit, head(collect({member: member, relationship: l})) AS chosen
WITH r, can_edit, chosen.member AS member, chosen.relationship AS relationship
ORDER BY id(r)
RETURN collect({
  recipe: r,
  can_edit: can_edit,
  added_by: {
    id: id(member),
    first_name: member.first_name,
    middle_name: member.middle_name,
    last_name: member.last_name,
    born: member.born,
    biological_sex: member.biological_sex,
    died: member.died,
    profile_picture: member.profile_picture
  },
  relationship: relationship
}) as entries
