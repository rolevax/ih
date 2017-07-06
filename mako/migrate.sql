create table replay_of_user(
	user_id bigint primary key references users(user_id),
	replay_ids bigint[] not null
);

insert into replay_of_user
(select t.user_id, array_agg(t.replay_id)
 from (select * from user_replay order by replay_id desc) as t
 group by t.user_id
);

drop table user_replay;

alter table user_girl drop column play;

alter table users alter level type smallint;
alter table users alter pt type smallint;
alter table girls alter level type smallint;
alter table girls alter pt type smallint;
alter table girls alter girl_id type integer;
alter table user_girl alter girl_id type integer;

alter table user_girl alter a_top type integer;
alter table user_girl alter a_last type integer;
alter table user_girl alter round type integer;
alter table user_girl alter win type integer;
alter table user_girl alter gun type integer;
alter table user_girl alter bark type integer;
alter table user_girl alter riichi type integer;
alter table user_girl alter ready type integer;

alter table user_girl alter a_top set default 0;
alter table user_girl alter a_last set default 0;
alter table user_girl alter round set default 0;
alter table user_girl alter win set default 0;
alter table user_girl alter gun set default 0;
alter table user_girl alter bark set default 0;
alter table user_girl alter riichi set default 0;
alter table user_girl alter ready set default 0;

alter table user_girl alter yaku_rci type integer; 
alter table user_girl alter yaku_ipt type integer;
alter table user_girl alter yaku_tmo type integer;
alter table user_girl alter yaku_tny type integer;
alter table user_girl alter yaku_pnf type integer;
alter table user_girl alter yaku_y1y type integer;
alter table user_girl alter yaku_y2y type integer;
alter table user_girl alter yaku_y3y type integer;
alter table user_girl alter yaku_jk1 type integer;
alter table user_girl alter yaku_jk2 type integer;
alter table user_girl alter yaku_jk3 type integer;
alter table user_girl alter yaku_jk4 type integer;
alter table user_girl alter yaku_bk1 type integer;
alter table user_girl alter yaku_bk2 type integer;
alter table user_girl alter yaku_bk3 type integer;
alter table user_girl alter yaku_bk4 type integer;
alter table user_girl alter yaku_ipk type integer;
alter table user_girl alter yaku_rns type integer;
alter table user_girl alter yaku_hai type integer;
alter table user_girl alter yaku_hou type integer;
alter table user_girl alter yaku_ckn type integer;
alter table user_girl alter yaku_ss1 type integer;
alter table user_girl alter yaku_it1 type integer;
alter table user_girl alter yaku_ct1 type integer;
alter table user_girl alter yaku_wri type integer;
alter table user_girl alter yaku_ss2 type integer;
alter table user_girl alter yaku_it2 type integer;
alter table user_girl alter yaku_ct2 type integer;
alter table user_girl alter yaku_toi type integer;
alter table user_girl alter yaku_ctt type integer;
alter table user_girl alter yaku_sak type integer;
alter table user_girl alter yaku_skt type integer;
alter table user_girl alter yaku_stk type integer;
alter table user_girl alter yaku_hrt type integer;
alter table user_girl alter yaku_s3g type integer;
alter table user_girl alter yaku_h1t type integer;
alter table user_girl alter yaku_jc2 type integer;
alter table user_girl alter yaku_mnh type integer;
alter table user_girl alter yaku_jc3 type integer;
alter table user_girl alter yaku_rpk type integer;
alter table user_girl alter yaku_c1t type integer;
alter table user_girl alter yaku_mnc type integer;
alter table user_girl alter yaku_x13 type integer;
alter table user_girl alter yaku_xd3 type integer;
alter table user_girl alter yaku_x4a type integer;
alter table user_girl alter yaku_xt1 type integer;
alter table user_girl alter yaku_xs4 type integer;
alter table user_girl alter yaku_xd4 type integer;
alter table user_girl alter yaku_xcr type integer;
alter table user_girl alter yaku_xr1 type integer;
alter table user_girl alter yaku_xth type integer;
alter table user_girl alter yaku_xch type integer;
alter table user_girl alter yaku_x4k type integer;
alter table user_girl alter yaku_x9r type integer;
alter table user_girl alter yaku_w13 type integer;
alter table user_girl alter yaku_w4a type integer;
alter table user_girl alter yaku_w9r type integer;
alter table user_girl alter kzeykm type integer;
alter table user_girl alter yaku_dora type integer;
alter table user_girl alter yaku_uradora type integer;
alter table user_girl alter yaku_akadora type integer;
alter table user_girl alter yaku_kandora  type integer;
alter table user_girl alter yaku_kanuradora type integer;

alter table user_girl alter yaku_rci set default 0;
alter table user_girl alter yaku_ipt set default 0;
alter table user_girl alter yaku_tmo set default 0;
alter table user_girl alter yaku_tny set default 0;
alter table user_girl alter yaku_pnf set default 0;
alter table user_girl alter yaku_y1y set default 0;
alter table user_girl alter yaku_y2y set default 0;
alter table user_girl alter yaku_y3y set default 0;
alter table user_girl alter yaku_jk1 set default 0;
alter table user_girl alter yaku_jk2 set default 0;
alter table user_girl alter yaku_jk3 set default 0;
alter table user_girl alter yaku_jk4 set default 0;
alter table user_girl alter yaku_bk1 set default 0;
alter table user_girl alter yaku_bk2 set default 0;
alter table user_girl alter yaku_bk3 set default 0;
alter table user_girl alter yaku_bk4 set default 0;
alter table user_girl alter yaku_ipk set default 0;
alter table user_girl alter yaku_rns set default 0;
alter table user_girl alter yaku_hai set default 0;
alter table user_girl alter yaku_hou set default 0;
alter table user_girl alter yaku_ckn set default 0;
alter table user_girl alter yaku_ss1 set default 0;
alter table user_girl alter yaku_it1 set default 0;
alter table user_girl alter yaku_ct1 set default 0;
alter table user_girl alter yaku_wri set default 0;
alter table user_girl alter yaku_ss2 set default 0;
alter table user_girl alter yaku_it2 set default 0;
alter table user_girl alter yaku_ct2 set default 0;
alter table user_girl alter yaku_toi set default 0;
alter table user_girl alter yaku_ctt set default 0;
alter table user_girl alter yaku_sak set default 0;
alter table user_girl alter yaku_skt set default 0;
alter table user_girl alter yaku_stk set default 0;
alter table user_girl alter yaku_hrt set default 0;
alter table user_girl alter yaku_s3g set default 0;
alter table user_girl alter yaku_h1t set default 0;
alter table user_girl alter yaku_jc2 set default 0;
alter table user_girl alter yaku_mnh set default 0;
alter table user_girl alter yaku_jc3 set default 0;
alter table user_girl alter yaku_rpk set default 0;
alter table user_girl alter yaku_c1t set default 0;
alter table user_girl alter yaku_mnc set default 0;
alter table user_girl alter yaku_x13 set default 0;
alter table user_girl alter yaku_xd3 set default 0;
alter table user_girl alter yaku_x4a set default 0;
alter table user_girl alter yaku_xt1 set default 0;
alter table user_girl alter yaku_xs4 set default 0;
alter table user_girl alter yaku_xd4 set default 0;
alter table user_girl alter yaku_xcr set default 0;
alter table user_girl alter yaku_xr1 set default 0;
alter table user_girl alter yaku_xth set default 0;
alter table user_girl alter yaku_xch set default 0;
alter table user_girl alter yaku_x4k set default 0;
alter table user_girl alter yaku_x9r set default 0;
alter table user_girl alter yaku_w13 set default 0;
alter table user_girl alter yaku_w4a set default 0;
alter table user_girl alter yaku_w9r set default 0;
alter table user_girl alter kzeykm set default 0;
alter table user_girl alter yaku_dora set default 0;
alter table user_girl alter yaku_uradora set default 0;
alter table user_girl alter yaku_akadora set default 0;
alter table user_girl alter yaku_kandora set default 0;
alter table user_girl alter yaku_kanuradora set default 0;

alter table user_girl add column ranks integer[4] not null default ARRAY[0,0,0,0];
alter table user_girl drop rank1;
alter table user_girl drop rank2;
alter table user_girl drop rank3;
alter table user_girl drop rank4;
update user_girl set ranks=ranks[0:3];

create function play(ranks integer[]) returns integer as
$$
BEGIN
  RETURN ranks[1]+ranks[2]+ranks[3]+ranks[4];
END;
$$ language plpgsql;

