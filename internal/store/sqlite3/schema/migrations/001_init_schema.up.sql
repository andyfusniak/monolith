begin immediate;

-- strftime('%Y-%m-%d %H:%M:%f000000+00:00', 'now')   <- current timestamp Go style

create table users (
  user_id       text primary key,
  email         text not null,
  password_hash text not null,
  created_at    text not null,
  constraint users_email_ukey unique (email)
) strict;

commit;
