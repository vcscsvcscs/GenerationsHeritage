MATCH (original:Recipe)<-[v:VariationOf]-(variation:Recipe)<-[:Created]-(creator:Person)
WHERE id(original) = $recipeId
RETURN variation, {
  notes: v.notes,
  created_at: v.created_at
} AS variation_relationship, {
  id: id(creator),
  first_name: creator.first_name,
  middle_name: creator.middle_name,
  last_name: creator.last_name,
  born: creator.born,
  biological_sex: creator.biological_sex,
  died: creator.died,
  profile_picture: creator.profile_picture
} AS creator
ORDER BY variation_relationship.created_at DESC
