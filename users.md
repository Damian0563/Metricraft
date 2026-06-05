# Supabase Users Table

Metricraft stores account and access-control data in the Supabase `public.users` table referenced by `DATABASE_USERS`.

## Schema

```sql
create table public.users (
  created_at timestamp with time zone not null default now(),
  app_name text null,
  mail text not null,
  secret text not null,
  uuid uuid null,
  allowed_users jsonb null default '[]'::jsonb,
  pending_users jsonb null default '[]'::jsonb,
  owner boolean null default false,
  constraint Users_pkey primary key (mail)
) TABLESPACE pg_default;
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `created_at` | `timestamp with time zone` | Creation timestamp, defaulting to `now()`. |
| `app_name` | `text` | Application or project name associated with the user. |
| `mail` | `text` | User email address. This is the primary key. |
| `secret` | `text` | Hashed sign-in secret stored by the backend. |
| `uuid` | `uuid` | Session/account token returned after sign-in. |
| `allowed_users` | `jsonb` | JSON array of email addresses allowed to access the owner's app. |
| `pending_users` | `jsonb` | JSON array of email addresses waiting for owner approval. |
| `owner` | `boolean` | Marks whether the row owns the app namespace. |
