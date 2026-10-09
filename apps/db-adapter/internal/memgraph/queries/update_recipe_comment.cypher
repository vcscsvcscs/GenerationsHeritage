MATCH (p:Person)-[c:CommentedOnRecipe]->(r:Recipe)
WHERE id(p) = $personId AND id(r) = $recipeId
SET c.message = $message,
    c.edited = $edited
RETURN {
  message: c.message,
  sent_at: c.sent_at,
  edited: c.edited
} AS comment, {
  id: id(p),
  first_name: p.first_name,
  middle_name: p.middle_name,
  last_name: p.last_name,
  born: p.born,
  biological_sex: p.biological_sex,
  died: p.died,
  profile_picture: p.profile_picture
} AS commenter
