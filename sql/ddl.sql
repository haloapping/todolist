CREATE TABLE
	users (
		id UUID UNIQUE NOT NULL,
		name VARCHAR NOT NULL,
		email VARCHAR UNIQUE NOT NULL,
		password VARCHAR(72) UNIQUE NOT NULL
	);

CREATE TABLE
	todos (
		id UUID UNIQUE NOT NULL,
		user_id UUID NOT NULL,
		title VARCHAR NOT NULL,
		description TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);