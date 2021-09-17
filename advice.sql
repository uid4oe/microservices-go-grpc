create table if not exists advices (
  user_id text not null,
  advice text not null,
  created_at timestamp,
  primary key (user_id)
);