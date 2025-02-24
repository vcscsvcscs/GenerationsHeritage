import type { Session, QueryResult, Transaction} from 'neo4j-driver';
import type{ Person,PersonProperties, FamilyTree } from '$lib/model';
import CreatePersonQuery from '$lib/server/queries/createPerson.cypher?raw';
import UpdatePersonQuery from '$lib/server/queries/updatePerson.cypher?raw';

export function createUser(db: Session, Person: PersonProperties): Promise<QueryResult<Person>> {
	return db.executeWrite(tx => tx.run<Person>(
		CreatePersonQuery, {Person})
	);
}

export function updateUser(db: Session, Person: PersonProperties): Promise<QueryResult<Person>> {
	return db.executeWrite(tx => tx.run<Person>(
		CreatePersonQuery, {Person})
	);
}

export function getUserFromGoogleId(db: Session,googleId: string): Promise<QueryResult<Person>> {
	return db.executeRead(tx => tx.run<Person>(""))
}

export interface User {
	id: number;
	email: string;
	googleId: string;
	name: string;
	picture: string;
}
