package role

const (
	createRole = `INSERT INTO role (name, type, rank, created_by, updated_by)
	VALUES (:name, :type, :rank, :created_by, :updated_by)`

	readRole = `
		SELECT
			id,
			name,
			type,
			rank,
			status,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM
			role`

	updateRole = `
	UPDATE
		role`

	readRoleCount = `
	SELECT
			COUNT(*)
		FROM
			role`
)
