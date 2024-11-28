CREATE TABLE IF NOT EXISTS mails (
		id SERIAL PRIMARY KEY,
		human_id INT NOT NULL,
		mail VARCHAR(255) NOT NULL,
        description VARCHAR(1024)
		);