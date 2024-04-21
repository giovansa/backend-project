package repository

var (
	qInsertUser = `
		INSERT INTO users(id, phone, name, password) 
		VALUES ($1, $2, $3, $4);`

	qGetUserByPhone = `
		SELECT
		    id,
		    phone,
		    name,
		    password,
		    created_at,
		    updated_at
		FROM users
		WHERE phone = $1;`

	qIncrementLoginCount = `
		UPDATE users
		SET success_login = success_login + 1,
		    updated_at = now()
		WHERE phone = $1;`
)
