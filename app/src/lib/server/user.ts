import type { Session } from 'neo4j-driver';

export function createUser(db: Session, googleId: string, email: string, first_name: string, family_name: string, middle_name: string): User {
	const row = db.run
	if (row === null) {
		throw new Error("Unexpected error");
	}
	const user: User = {
		id: row.number(0),
		googleId,
		email,
	};
	return user;
}

export function getUserFromGoogleId(googleId: string): User | null {
	const row = db.queryOne("SELECT id, google_id, email, name, picture FROM user WHERE google_id = ?", [googleId]);
	if (row === null) {
		return null;
	}
	const user: User = {
		id: row.number(0),
		googleId: row.string(1),
		email: row.string(2),
		name: row.string(3),
		picture: row.string(4)
	};
	return user;
}

export interface User {
	id: number;
	email: string;
	googleId: string;
	name: string;
	picture: string;
}
