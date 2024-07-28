package url

const (
	createUrl = `
	INSERT INTO url (original_url, shorten_url, fk_user_id, created_by, updated_by)
	    VALUES (:original_url, :shorten_url, :fk_user_id, :created_by, :updated_by)`

	readUrl = `
	SELECT
	    id,
		original_url,
		shorten_url,
	    visit,
	    fk_user_id,
	    status,
	    created_at,
	    created_by,
	    updated_at,
	    updated_by
	FROM
	    url`

	updateUrl = `
	UPDATE
	    url`

	readUrlCount = `
	SELECT
	    COUNT(*)
	FROM
	    url`
)
