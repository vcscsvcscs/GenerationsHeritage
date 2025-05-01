CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.last_name);
CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.first_name);
CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.born);
CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.mothers_first_name);
CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.mothers_last_name);
CREATE CONSTRAINT ON (n:Person) ASSERT n.google_id IS UNIQUE;
CREATE CONSTRAINT ON (n:Person) ASSERT n.last_name, n.first_name, n.born, n.mothers_first_name, n.mothers_last_name IS UNIQUE;
