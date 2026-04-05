MATCH (original:Recipe)<-[v:VariationOf]-(variation:Recipe)<-[:Created]-(creator:Person)
WHERE id(original) = $recipeId
RETURN variation, v, creator
ORDER BY v.created_at DESC
