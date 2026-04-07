MATCH (p:Person)-[c:CommentedOnRecipe]->(r:Recipe)
WHERE id(r) = $recipeId
RETURN collect({
    comment: c,
    commenter: {
        id: id(p),
        first_name: p.first_name,
        last_name: p.last_name,
        profile_picture: p.profile_picture
    }
}) as comments
