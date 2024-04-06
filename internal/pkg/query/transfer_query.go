package query

const (
	CREATE_NEW_TRANSFER_WITH_RETURNING_QUERY = `
		INSERT INTO 
			transfers (recipient_account_id, sender_account_id, amount, status) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING 
			id, recipient_account_id, sender_account_id, amount, status, created_at;
	`

	UPDATE_TRANSFER_STATUS_BY_ID = `
		UPDATE 
			transfers 
		SET 
			status = $2, completed_at = NOW()
		WHERE 
			id = ($1);`
)
