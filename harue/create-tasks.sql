CREATE TABLE tasks (
	task_id int CONSTRAINT tasks_pk PRIMARY KEY,
	title varchar(64) NOT NULL,
	content text NOT NULL,
	state integer NOT NULL DEFAUlT 0,
	assignee_id bigint REFERENCES users(user_id),
	c_point int NOT NULL DEFAULT 0
);

