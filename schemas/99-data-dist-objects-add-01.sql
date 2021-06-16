insert into links(id, hash, published, usage_limit, created_by, created_at, expired_at)
values (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day),
       (nextval('links_id_seq'), encode(crypt(concat_ws('data:', nextval('links_helper_seq'), now()), gen_salt('des'))::bytea, 'base64'), true, 10, 1, now(), now() + interval '3' day)
;

insert into distribution_objects(id, distribution_id, recipient_id, respondent_id, link_id, created_by, created_at)
values
    /** rifqi: samuel, joseph, jocelyn, yakub, coktra */
    (nextval('dist_object_id_seq'), 1, 25, 80, 256, 1, now()),
    (nextval('dist_object_id_seq'), 1, 25, 81, 257, 1, now()),
    (nextval('dist_object_id_seq'), 1, 25, 82, 258, 1, now()),
    (nextval('dist_object_id_seq'), 1, 25, 55, 259, 1, now()),
    (nextval('dist_object_id_seq'), 1, 25, 47, 260, 1, now()),

    /** melinda: leviana, arindi, dhony, indra, nindyasuri */
    (nextval('dist_object_id_seq'), 1, 23, 75, 261, 1, now()),
    (nextval('dist_object_id_seq'), 1, 23, 76, 262, 1, now()),
    (nextval('dist_object_id_seq'), 1, 23, 77, 263, 1, now()),
    (nextval('dist_object_id_seq'), 1, 23, 78, 264, 1, now()),
    (nextval('dist_object_id_seq'), 1, 23, 79, 265, 1, now())
;