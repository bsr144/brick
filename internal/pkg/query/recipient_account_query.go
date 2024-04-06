package query

const (
	// Write queries
	CREATE_NEW_RECIPIENT_ACCOUNT_WITH_RETURNING_QUERY = `
		INSERT INTO 
			recipient_accounts (account_number, account_name, bank_code, bank_name, verification_status) 
		VALUES 
			($1, $2, $3, $4, $5) 
		RETURNING 
			id, account_number, account_name, bank_code, bank_name, verification_status;
	`
	UPDATE_RECIPIENT_ACCOUNT_BY_ID = `
		UPDATE 
			recipient_accounts 
		SET 
			account_number = $2, account_name = $3, bank_code = $4, bank_name = $5, verification_status = $6, last_verified_at = NOW()
		WHERE 
			id = ($1)
		RETURNING
			id, account_number, account_name, bank_code, bank_name, verification_status, last_verified_at;`

	// Read queries
	GET_RECIPIENT_ACCOUNT_BY_BANK_CODE_AND_BANK_ACCOUNT_NUMBER = `
		SELECT * 
		FROM 
			recipient_accounts
		WHERE 
			bank_code = ($1) AND account_number = ($2);`
)
