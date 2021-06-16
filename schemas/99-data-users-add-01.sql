
insert into users(id, username, display_name, email, phone, password, access_role, job_role, assignment, group_id, created_at)
values
    /* YB peers */
    (nextval('users_id_seq'), 'leviana.hestiawan', 'Leviana Hestiawan', 'leviana.hestiawan@gmail.com', null, null, 'member', 'Digital Active Squad', 'Digital Active Squad', 2, now()),
    (nextval('users_id_seq'), 'arindiassari23', 'Arindi Assari', 'arindiassari23@gmail.com', null, null, 'member', 'Digital Active Squad', 'Digital Active Squad', 2, now()),
    (nextval('users_id_seq'), 'dhonyrulan', 'Dhony Rulan', 'dhonyrulan@gmail.com', null, null, 'member', 'Digital Active Squad', 'Digital Active Squad', 2, now()),
    (nextval('users_id_seq'), 'indra.rukasyah', 'Indra Rukasyah', 'indra.rukasyah@gmail.com', null, null, 'member', 'Digital Active Squad', 'Digital Active Squad', 2, now()),
    (nextval('users_id_seq'), 'f.nindyasuri', 'F. Nindyasuri', 'f.nindyasuri@gmail.com', null, null, 'member', 'Digital Active Squad', 'Digital Active Squad', 2, now()),

    (nextval('users_id_seq'), 'samuel.christian7', 'Samuel Christian Hidajat', 'samuel.christian7@gmail.com', null, null, 'member', 'YB', 'Propo dan Campaign Engine', 2, now()),
    (nextval('users_id_seq'), 'joseph', 'Joseph', 'joseph.ocbcnisp@gmail.com', null, null, 'member', 'Software Engineer', 'Architecture And Digital', 2, now()),
    (nextval('users_id_seq'), 'jocelyntjahyadi26', 'Jocelyn Olivia Tjahyadi', 'jocelyntjahyadi26@gmail.com', null, null, 'member', 'Software Engineer', 'Architecture And Digital', 2, now())
;