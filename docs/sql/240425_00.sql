-- [DML] Populate role admin and user
INSERT INTO `role` (`name`, `type`, `rank`) VALUES
('Super Admin', 'admin', 1),
('User', 'user', 2)
;

-- [DMLI nsertion on Action Table
INSERT INTO `action`
(`id`, `code`, `description`)
VALUES
(1, 'api.admin.list', 'Get Listing API that only for Admin')
;

-- [DML] Insertion on User Table