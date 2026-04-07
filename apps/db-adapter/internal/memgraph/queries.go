package memgraph

import _ "embed"

//go:embed queries/create_indexes.cypher
var CreateIndexesCypherQuery string

//go:embed queries/drop_indexes.cypher
var DropIndexesCypherQuery string

//go:embed queries/create_constraints.cypher
var CreateConstraintsCypherQuery string

//go:embed queries/drop_constraints.cypher
var DropConstraintsCypherQuery string

// Requires Person parameter.
// Returns p as person
//
//go:embed queries/create_person.cypher
var CreatePersonCypherQuery string

// Requires id parameter.
// Returns n as person
//
//go:embed queries/get_person_by_id.cypher
var GetPersonCypherQuery string

// Requires id, props parameter.
// Returns n as person
//
//go:embed queries/update_person.cypher
var UpdatePersonCypherQuery string

// Requires invite_code, props parameter.
//
//go:embed queries/update_person_by_invite_code.cypher
var UpdatePersonByInviteCodeCypherQuery string

// Requires google_id parameter.
// Returns n as person
//
//go:embed queries/get_person_by_google_id.cypher
var GetPersonByGoogleIdCypherQuery string

// Requires id parameter.
//
// Returns labels(n) AS labels, n AS person
//
//go:embed queries/soft_delete_person_by_id.cypher
var SoftDeletePersonCypherQuery string

// Requires id parameter.
//
//go:embed queries/hard_delete_person_by_id.cypher
var HardDeletePersonCypherQuery string

// Requires id1, id2 parameters.
//
// returns relationship
//
//go:embed queries/get_relationship.cypher
var GetRelationshipCypherQuery string

// Requires childId, parentId, childRelationship, parentRelationship parameters.
//
// returns relationships
//
//go:embed queries/create_child_parent_relationships.cypher
var CreateChildParentRelationshipCypherQuery string

// Requires childId, parentId parameters.
//
// returns relationships
//
//go:embed queries/create_sibling_relationships_based_on_parent.cypher
var CreateSiblingRelationshipsBasedOnParentCypherQuery string

// Requires id1, id2, Relationship1, Relationship1 parameters.
//
// return relationships
//
//go:embed queries/create_sibling_relationship.cypher
var CreateSiblingRelationshipCypherQuery string

// Requires id1, id2, Relationship1, Relationship1 parameters.
//
// return relationships
//
//go:embed queries/create_spouse_relationship.cypher
var CreateSpouseRelationshipCypherQuery string

// Requires id1, id2 parameters.
//
//go:embed queries/delete_relationship.cypher
var DeleteRelationshipCypherQuery string

// Requires id1, id2, relationship parameter.
//
// return relationship
//
//go:embed queries/update_relationship.cypher
var UpdateRelationshipCypherQuery string

// Requires id1, id2 parameter.
//
// return relationship
//
//go:embed queries/create_admin_relationship.cypher
var CreateAdminRelationshipCypherQuery string

// Requires id1, id2 parameter.
//
//go:embed queries/delete_admin_relationship.cypher
var DeleteAdminRelationshipCypherQuery string

// Requires id1, id2 parameter.
//
// return relationship
//
//go:embed queries/get_admin_relationship.cypher
var GetAdminRelationshipCypherQuery string

// Requires id parameter.
//
// returns admins
//
//go:embed queries/get_profile_admins.cypher
var GetProfileAdminsCypherQuery string

// Requires id parameter.
//
// returns managed
//
//go:embed queries/get_managed_profiles.cypher
var GetManagedProfilesCypherQuery string

// Requires id parameter.
//
// returns people, relationships
//
//go:embed queries/get_blood_relations_by_id.cypher
var GetBloodRelativesCypherQuery string

// Requires id parameter.
//
// returns people, relationships
//
//go:embed queries/get_family_tree_with_spouses.cypher
var GetFamilyTreeWithSpousesCypherQuery string

// Requires comment, id1 as commenter and id2 as profile that is commented on parameter.
//
// returns people, comments
//
//go:embed queries/comment.cypher
var CommentCypherQuery string

// Requires id1 as commenter and id2 as profile that is commented on parameter.
//
//go:embed queries/delete_comment.cypher
var DeleteCommentCypherQuery string

// Requires id1 as profile that is commented on parameter.
//
// returns comments, people
//
//go:embed queries/comments_on_profile.cypher
var CommentsOnProfileCypherQuery string

// Requires personId, Recipe, Relationship parameters.
//
// returns recipe, relationship
//
//go:embed queries/create_recipe_with_relationship.cypher
var CreateRecipeWithRelationshipCypherQuery string

// Requires id parameter.
//
// returns recipes, recipeRelations
//
//go:embed queries/get_recipes_by_person_id.cypher
var GetRecipesByPersonIdCypherQuery string

// Requires id, props parameters.
//
// returns recipe
//
//go:embed queries/update_recipe.cypher
var UpdateRecipeCypherQuery string

// Requires id parameter.
//
// returns labels, recipe
//
//go:embed queries/soft_delete_recipe.cypher
var SoftDeleteRecipeCypherQuery string

// Requires id parameter.
//
//go:embed queries/hard_delete_recipe.cypher
var HardDeleteRecipeCypherQuery string

// Requires personId, recipeId, Relationship parameters.
//
// returns relationship
//
//go:embed queries/create_recipe_relationship.cypher
var CreateRecipeRelationshipCypherQuery string

// Requires personId, recipeId parameters.
//
//go:embed queries/delete_recipe_relationship.cypher
var DeleteRecipeRelationshipCypherQuery string

// Requires recipeId, userId parameters.
// Returns r (recipe) if user is allowed to manage it.
//
//go:embed queries/could_manage_recipe.cypher
var CouldManageRecipeCypherQuery string

// Requires id parameter. Distance is embedded via fmt.Sprintf.
//
// returns entries
//
//go:embed queries/get_family_cookbook.cypher
var GetFamilyCookbookCypherQueryTemplate string

// Requires recipeId, userId parameters.
// Returns r (recipe) if user is allowed to see it.
//
//go:embed queries/could_see_recipe.cypher
var CouldSeeRecipeCypherQuery string

// Requires originalRecipeId, creatorId, RecipeProperties, VariationProperties, LikesProperties parameters.
//
// returns recipe, variation_relationship, likes_relationship
//
//go:embed queries/create_recipe_variation.cypher
var CreateRecipeVariationCypherQuery string

// Requires recipeId parameter.
//
// returns variations, variation_relationships, creators
//
//go:embed queries/get_recipe_variations.cypher
var GetRecipeVariationsCypherQuery string

// Requires personId, recipeId, Comment parameters.
//
// returns comment, commenter
//
//go:embed queries/comment_on_recipe.cypher
var CommentOnRecipeCypherQuery string

// Requires recipeId parameter.
//
// returns comments (array with commenter info)
//
//go:embed queries/get_recipe_comments.cypher
var GetRecipeCommentsCypherQuery string

// Requires personId, recipeId, message parameters.
//
// returns comment
//
//go:embed queries/update_recipe_comment.cypher
var UpdateRecipeCommentCypherQuery string

// Requires personId, recipeId parameters.
//
//go:embed queries/delete_recipe_comment.cypher
var DeleteRecipeCommentCypherQuery string
