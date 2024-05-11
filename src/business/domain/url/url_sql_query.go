package url

const (
	createUrl = `
	INSERT INTO role (original_url, shorten_url, user_id, created_by, updated_by)
	    VALUES (:original_url, :shorten_url, :user_id, :created_by, :updated_by)`

	getUrl = `
	SELECT
	    id,
		original_url,
		shorten_url,
	    visit,
	    user_id,
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
