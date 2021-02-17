package model

type User struct {
	ID       int     `json:"id"`
	Account  string  `json:"account"`
	Name     *string `json:"name"`
	Password string  `json:"password"`
}

const UserSQL = `
CREATE TABLE IF NOT EXISTS users (
	id serial primary key,
	account varchar(60) unique not null,
	name varchar(255) null,
	password varchar(255) not null
);
`
