package query

const (
	// Write queries
	CREATE_NEW_CREDENTIAL = `
		INSERT INTO 
			api_credentials (client_id, client_secret, user_id) 
		VALUES ($1, $2, $3);
	`

	// Read queries
	GET_CREDENTIAL_BY_USER_ID = `
	SELECT 
		id, client_id, client_secret, user_id 
	FROM 
		api_credentials 
	WHERE 
		user_id = ($1);`

	GET_CREDENTIAL_BY_CLIENT_ID = `
	SELECT 
		id, client_id, client_secret, user_id 
	FROM 
		api_credentials 
	WHERE 
		client_id = ($1);`
)
