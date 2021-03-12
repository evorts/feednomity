/* sample distributions */
-- noinspection SqlWithoutWhere
delete from links;
-- noinspection SqlWithoutWhere
delete from distribution_objects;
-- noinspection SqlWithoutWhere
delete from distributions;

insert into distributions(id, topic, distribution_limit, distribution_count, created_at)
values
(1, 'Q1+Q2 360 Review - Beyond Banking', 3, 0, now())
;

insert into distribution_objects(id, distribution_id, recipient_id, respondent_id, created_at)
values
(1, 1, 36, 36, now())
;

insert into links(id, distribution_object_id, hash, published, usage_limit, created_at)
values
(1, 1, crypt(concat_ws('data:',1,1,0, now()), gen_salt('des')), true, 0, now())
;