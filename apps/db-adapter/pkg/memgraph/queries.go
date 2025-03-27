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
//go:embed queries/create_person.cypher
var CreatePersonCypherQuery string

// Requires id parameter.
// Returns n as person
//go:embed queries/get_person_by_id.cypher
var GetPersonCypherQuery string

// Requires id, props parameter.
// Returns n as person
//go:embed queries/update_person.cypher
var UpdatePersonCypherQuery string

// Requires google_id parameter.
// Returns n as person
//go:embed queries/get_person_by_google_id.cypher
var GetPersonByGoogleIdCypherQuery string

// Requires id parameter.
// Returns labels(n) AS labels, n AS person
//go:embed queries/soft_delete_person_by_id.cypher
var SoftDeletePersonCypherQuery string

// Requires id parameter.
//go:embed queries/hard_delete_person_by_id.cypher
var HardDeletePersonCypherQuery string

// Requires id1, id2 parameters.
//go:embed queries/get_relationship.cypher
var GetRelationshipCypherQuery string

// Requires id1, id2, Relationship parameters.
//go:embed queries/create_directed_relationship.cypher
var CreateDirectedRelationshipCypherQuery string

// Requires id1, id2, Relationship1, Relationship2 parameters.
//go:embed queries/create_two_directed_relationships.cypher
var CreateTwoDirectedRelationshipCypherQuery string

// Requires id parameter.
//go:embed queries/get_family_tree_by_id.cypher
var GetFamilyTreeByIdCypherQuery string
