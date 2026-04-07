MATCH (original:Recipe), (creator:Person)
WHERE id(original) = $originalRecipeId AND id(creator) = $creatorId
CREATE (variation:Recipe $RecipeProperties)
CREATE (creator)-[:Created]->(variation)
CREATE (variation)-[v:VariationOf $VariationProperties]->(original)
CREATE (creator)-[l:Likes $LikesProperties]->(variation)
RETURN variation as recipe, v as variation_relationship, l as likes_relationship