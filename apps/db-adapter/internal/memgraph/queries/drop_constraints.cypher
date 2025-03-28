DROP CONSTRAINT ON (n:Person) ASSERT EXISTS (n.last_name);
DROP CONSTRAINT ON (n:Person) ASSERT EXISTS (n.first_name);
DROP CONSTRAINT ON (n:Person) ASSERT EXISTS (n.born);
DROP CONSTRAINT ON (n:Person) ASSERT EXISTS (n.mothers_first_name);
DROP CONSTRAINT ON (n:Person) ASSERT EXISTS (n.mothers_last_name);
DROP CONSTRAINT ON (n:Person) ASSERT n.google_id IS UNIQUE;
DROP CONSTRAINT ON (n:Person) ASSERT n.last_name, n.first_name, n.born, n.mothers_first_name, n.mothers_last_name IS UNIQUE;
