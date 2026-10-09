MATCH (original:Recipe), (creator:Person)
WHERE id(original) = $originalRecipeId AND id(creator) = $creatorId
CREATE (variation:Recipe $RecipeProperties)
CREATE (creator)-[:Created]->(variation)
CREATE (variation)-[v:VariationOf $VariationProperties]->(original)
CREATE (creator)-[l:Likes $LikesProperties]->(variation)
RETURN variation AS recipe, {
  notes: v.notes,
  created_at: v.created_at
} AS variation_relationship, l AS likes_relationship
