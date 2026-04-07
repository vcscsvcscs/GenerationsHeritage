MATCH (n:Person)
WHERE id(n) = $id
OPTIONAL MATCH (n)-[:Parent|Child|Sibling|Spouse*1..%d]-(relative:Person)
WITH collect(DISTINCT relative) + [n] AS family
UNWIND family AS member
OPTIONAL MATCH (member)-[l:Likes]->(r:Recipe)
WHERE r IS NOT NULL
RETURN collect(DISTINCT {
  recipe: r,
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
  relationship: l
}) as entries
