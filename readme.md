```postgresql
create table entities(
  id bigserial primary key,
  code text,
  value bigint
);

begin;
do $do$
  begin
    for i in 1..1000000 loop
      insert into entities(code, value) values ((select substr(md5(random()::text), 1, 25)), i*1000);
      end loop;

  end $do$;
commit;

select * from entities;
```