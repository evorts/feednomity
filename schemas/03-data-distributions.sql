/* sample distributions */
-- noinspection SqlWithoutWhere
delete
from distribution_objects;
-- noinspection SqlWithoutWhere
delete
from distributions;

insert into distributions(id, topic, distribution_limit, distribution_count, range_start, range_end, created_by,
                          for_group_id, created_at)

values (1, '360 Review - Beyond Banking (Des 2020 ~ May 2021)', 5, 0, to_date('20201201', 'YYYYMMDD'),
        to_date('20210531', 'YYYYMMDD'), 1, 2, now()),
       (2, '360 Review - Beyond Banking (Jun 2021 ~ Nov 2021)', 5, 0, to_date('20210601', 'YYYYMMDD'),
        to_date('20211130', 'YYYYMMDD'), 1, 2, now())
;

alter sequence if exists links_helper_seq restart with 1;
insert into distribution_objects(id, distribution_id, recipient_id, respondent_id, link_id, created_by, created_at)
values
        /** abraham tobing: tomo, silvia, okta, didi, kipli */
        (nextval('dist_object_id_seq'), 1, 6, 66, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 6, 36, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 6, 30, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 6, 16, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 6, 29, nextval('links_helper_seq'), 1, now()),

        /* iman: julius, yeka, ian, rianita, tiolie, didi, ferdy, ervina, safandhi, silvia  */
        (nextval('dist_object_id_seq'), 1, 7, 65, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 55, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 21, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 32, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 38, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 16, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 27, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 34, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 7, 36, nextval('links_helper_seq'), 1, now()),

        /* aldo: fenny, denza, haris, wibi, diana, alan  */
        (nextval('dist_object_id_seq'), 1, 9, 63, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 9, 53, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 9, 20, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 9, 39, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 9, 15, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 9, 57, nextval('links_helper_seq'), 1, now()),

        /* irul: yvonne, yuni, nicko, ervina, ferdy, tiolie  */
        (nextval('dist_object_id_seq'), 1, 10, 70, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 10, 41, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 10, 26, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 10, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 10, 27, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 10, 38, nextval('links_helper_seq'), 1, now()),

        /* angga: lucy, kipli, calista, ozi, samuel  */
        (nextval('dist_object_id_seq'), 1, 11, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 11, 29, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 11, 13, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 11, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 11, 35, nextval('links_helper_seq'), 1, now()),

        /* bistok: resky, lucy, jhoni, samuel, edwin, ruminta, kipli  */
        (nextval('dist_object_id_seq'), 1, 12, 71, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 12, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 12, 64, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 12, 35, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 12, 18, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 12, 33, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 12, 29, nextval('links_helper_seq'), 1, now()),

        /* calista: lucy, samuel, ozi, denza, angga, bistok, santi */
        (nextval('dist_object_id_seq'), 1, 13, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 13, 35, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 13, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 13, 53, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 13, 11, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 13, 12, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 13, 44, nextval('links_helper_seq'), 1, now()),

        /* dede: resky, lucy, edwin, ozi, novita, santi */
        (nextval('dist_object_id_seq'), 1, 14, 71, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 14, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 14, 18, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 14, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 14, 28, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 14, 44, nextval('links_helper_seq'), 1, now()),

        /* diana: fenny, wibi, alan, aldo, haris, ozi */
        (nextval('dist_object_id_seq'), 1, 15, 63, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 15, 39, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 15, 57, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 15, 9, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 15, 20, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 15, 31, nextval('links_helper_seq'), 1, now()),

        /* didi: purwandi, coktra, bistok, edwin, tiolie, nicko, kipli, bowo, iman, wibi */
        (nextval('dist_object_id_seq'), 1, 16, 52, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 47, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 12, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 18, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 38, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 26, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 29, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 54, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 7, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 16, 39, nextval('links_helper_seq'), 1, now()),

        /* echo: fenny, wibi, diana, alan, bowo */
        (nextval('dist_object_id_seq'), 1, 17, 63, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 17, 39, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 17, 15, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 17, 57, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 17, 54, nextval('links_helper_seq'), 1, now()),

        /* edwin: resky, lucy, samuel, bistok, dede, kipli, novita */
        (nextval('dist_object_id_seq'), 1, 18, 71, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 18, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 18, 35, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 18, 12, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 18, 14, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 18, 29, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 18, 28, nextval('links_helper_seq'), 1, now()),

        /* ervina: yvonne, iman, irul, ade, ozi, ferdy, yuni */
        (nextval('dist_object_id_seq'), 1, 19, 70, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 19, 7, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 19, 10, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 19, 59, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 19, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 19, 27, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 19, 41, nextval('links_helper_seq'), 1, now()),

        /* haris: fenny, alan, wibi, kipli, diana, denza */
        (nextval('dist_object_id_seq'), 1, 20, 63, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 20, 57, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 20, 39, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 20, 29, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 20, 15, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 20, 53, nextval('links_helper_seq'), 1, now()),

        /* ian: julius, iman, rianita, bowo */
        (nextval('dist_object_id_seq'), 1, 21, 65, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 21, 7, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 21, 32, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 21, 54, nextval('links_helper_seq'), 1, now()),

        /* melinda: ? */

        /* ghozy: coktra, purwandi, komang */
        (nextval('dist_object_id_seq'), 1, 24, 47, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 24, 52, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 24, 73, nextval('links_helper_seq'), 1, now()),

        /* rifqi: ghozy */
        (nextval('dist_object_id_seq'), 1, 25, 24, nextval('links_helper_seq'), 1, now()),

        /* ferdy: yvonne, silvia, ervina, iman, safandhi */
        (nextval('dist_object_id_seq'), 1, 27, 70, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 27, 36, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 27, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 27, 7, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 27, 34, nextval('links_helper_seq'), 1, now()),

        /* novita: lucy, ozi, samuel, bowo, calista */
        (nextval('dist_object_id_seq'), 1, 28, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 28, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 28, 35, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 28, 54, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 28, 13, nextval('links_helper_seq'), 1, now()),

        /* kipli: resky, jhoni, angga, ruminta, dede, bistok */
        (nextval('dist_object_id_seq'), 1, 29, 71, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 29, 64, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 29, 11, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 29, 33, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 29, 14, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 29, 12, nextval('links_helper_seq'), 1, now()),

        /* okta: tomo, abe, silvia, ozi, vkc */
        (nextval('dist_object_id_seq'), 1, 30, 66, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 30, 6, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 30, 36, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 30, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 30, 69, nextval('links_helper_seq'), 1, now()),

        /* ozi: lucy, irwin, novita, samuel, calista, dede */
        (nextval('dist_object_id_seq'), 1, 31, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 31, 22, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 31, 28, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 31, 35, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 31, 13, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 31, 14, nextval('links_helper_seq'), 1, now()),

        /* rianita: julius, iman, tiolie, ozi, ian */
        (nextval('dist_object_id_seq'), 1, 32, 65, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 32, 7, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 32, 38, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 32, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 32, 21, nextval('links_helper_seq'), 1, now()),

        /* ruminta: jhoni, bistok, kipli, ozi, angga */
        (nextval('dist_object_id_seq'), 1, 33, 64, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 33, 12, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 33, 29, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 33, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 33, 11, nextval('links_helper_seq'), 1, now()),

        /* safandhi: yvonne, ade, ervina, bowo, ferdy, iman */
        (nextval('dist_object_id_seq'), 1, 34, 70, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 34, 59, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 34, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 34, 54, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 34, 27, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 34, 7, nextval('links_helper_seq'), 1, now()),

        /* samuel: lucy, edwin, calista, angga, novita, bistok */
        (nextval('dist_object_id_seq'), 1, 35, 67, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 35, 18, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 35, 13, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 35, 11, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 35, 28, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 35, 12, nextval('links_helper_seq'), 1, now()),

        /* silvia: yvonne, tomo, abe, ade, ferdy, irul */
        (nextval('dist_object_id_seq'), 1, 36, 70, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 36, 66, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 36, 6, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 36, 59, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 36, 27, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 36, 10, nextval('links_helper_seq'), 1, now()),

        /* tiolie: julius, rianita, irul, iman, ervina */
        (nextval('dist_object_id_seq'), 1, 38, 65, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 38, 32, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 38, 10, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 38, 7, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 38, 19, nextval('links_helper_seq'), 1, now()),

        /* wibi: fenny, nico, safandhi, aldo, haris, echo, didi, diana */
        (nextval('dist_object_id_seq'), 1, 39, 63, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 39, 72, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 39, 34, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 39, 9, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 39, 20, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 39, 17, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 39, 16, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 39, 15, nextval('links_helper_seq'), 1, now()),

        /** yuni evalin: yvonne, nicko prasetyo, ervina, silvia, ferdy, denza */
        (nextval('dist_object_id_seq'), 1, 41, 70, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 41, 26, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 41, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 41, 36, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 41, 27, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 41, 53, nextval('links_helper_seq'), 1, now()),

        /** impola: anita, santi, ozi, dede, ervina, diana, rianita */
        (nextval('dist_object_id_seq'), 1, 43, 42, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 43, 44, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 43, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 43, 14, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 43, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 43, 15, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 43, 32, nextval('links_helper_seq'), 1, now()),

        /** anita: impola, santi, ozi, dede, ervina, diana, rianita */
        (nextval('dist_object_id_seq'), 1, 42, 43, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 42, 44, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 42, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 42, 14, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 42, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 42, 15, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 42, 32, nextval('links_helper_seq'), 1, now()),

        /** santi: impola, anita, ozi, dede, ervina, diana, rianita */
        (nextval('dist_object_id_seq'), 1, 44, 43, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 44, 42, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 44, 31, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 44, 14, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 44, 19, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 44, 15, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 44, 32, nextval('links_helper_seq'), 1, now()),

        /** purwandi: didi, coktra, ghozy, wibi, iman, bistok */
        (nextval('dist_object_id_seq'), 1, 52, 16, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 52, 47, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 52, 24, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 52, 39, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 52, 7, nextval('links_helper_seq'), 1, now()),
        (nextval('dist_object_id_seq'), 1, 52, 12, nextval('links_helper_seq'), 1, now())
;

