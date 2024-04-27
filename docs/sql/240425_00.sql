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
INSERT INTO `user`
(`fk_role_id`, `email`, `username`, `password`, `display_name`)
VALUES
(0, `adiatma85@gmail.com`, `adiatma85`, `$2a$10$FnIk9LoiSBYw6gHTNmg4redG.9EKWCO5K3h0zE.4KTfWFj9aNfRhK`, 'Ramdani Koernia');
