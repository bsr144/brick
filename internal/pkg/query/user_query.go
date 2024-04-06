package query

const (
	// Write queries
	CREATE_NEW_USER_WITH_RETURNING_QUERY = `
		INSERT INTO 
			users (email, password, salt, balance) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING 
			id, email;
	`

	// Read queries
	GET_USER_BY_EMAIL_QUERY = `
		SELECT * 
		FROM 
			users 
		WHERE 
			email = ($1) AND deleted_at IS NULL;`

	GET_USER_BY_ID = `
		SELECT 
			id, email, balance 
		FROM 
			users 
		WHERE 
			id = ($1) AND deleted_at IS NULL;`

	GET_USER_INFO_BY_ID = `
	SELECT 
    	sa.id, sa.email, sa.balance, ac.id, ac.client_id, ac.client_secret, ac.sender_account_id
	FROM 
		users sa
	JOIN
		api_credentials ac ON sa.id = ac.sender_account_id
	WHERE
		sa.id = ($1) AND sa.deleted_at IS NULL;`
)
