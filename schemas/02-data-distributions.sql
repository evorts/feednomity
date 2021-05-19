/* sample distributions */
-- noinspection SqlWithoutWhere
delete from distribution_objects;
-- noinspection SqlWithoutWhere
delete from distributions;
-- noinspection SqlWithoutWhere
delete from links;

insert into links(id, hash, published, usage_limit, created_at)
values
(1, crypt(concat_ws('data:',1, now()), gen_salt('des')), true, 10, now()),
(2, crypt(concat_ws('data:',2, now()), gen_salt('des')), true, 10, now()),
(3, crypt(concat_ws('data:',3, now()), gen_salt('des')), true, 10, now()),
(4, crypt(concat_ws('data:',4, now()), gen_salt('des')), true, 10, now()),
(5, crypt(concat_ws('data:',5, now()), gen_salt('des')), true, 10, now())
;

insert into distributions(id, topic, distribution_limit, distribution_count, range_start, range_end, created_by, for_group_id, created_at)
values
(1, '360 Review - Beyond Banking (Des 2020 ~ May 2021)', 5, 0, to_date('20201201','YYYYMMDD'), to_date('20210531','YYYYMMDD'), 1, 2, now()),
(2, '360 Review - Beyond Banking (Jun 2021 ~ Nov 2021)', 5, 0, to_date('20210601','YYYYMMDD'), to_date('20211130','YYYYMMDD'), 1, 2, now())
;

insert into distribution_objects(id, distribution_id, recipient_id, respondent_id, link_id, created_by, created_at)
values
(1, 1, 2, 3, 1, 1, now()),
(2, 1, 2, 4, 2, 1, now()),
(3, 1, 2, 5, 3, 1, now())
;

