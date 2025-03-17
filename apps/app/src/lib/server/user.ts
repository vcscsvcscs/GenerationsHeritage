import type { Session, QueryResult } from 'neo4j-driver';
import type { Person, PersonProperties } from '$lib/model';
import CreatePersonQuery from '$lib/server/queries/create_person.cypher?raw';
import UpdatePersonQuery from '$lib/server/queries/update_person.cypher?raw';
import GetPersonByGoogleID from '$lib/server/queries/get_person_by_google_id.cypher?raw';

export function createUser(db: Session, Person: PersonProperties): Promise<QueryResult<Person>> {
	return db.executeWrite(tx => tx.run<Person>(
		CreatePersonQuery, { props: Person })
	);
}

export function updateUser(db: Session, Person: PersonProperties): Promise<QueryResult<Person>> {
	return db.executeWrite(tx => tx.run<Person>(
		UpdatePersonQuery, { props: Person })
	);
}

export function getUserFromGoogleId(db: Session, googleID: string): Promise<QueryResult<Person>> {
	return db.executeRead(tx => tx.run<Person>(
		GetPersonByGoogleID, { google_id: googleID })
	);
}
