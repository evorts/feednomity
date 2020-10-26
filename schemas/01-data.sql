/** insert default super admin **/
truncate users;
insert into users(username, display_name, role, created_date)
values ('sysadmin', 'System Administrator', 'sysadmin', now());