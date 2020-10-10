# Database IDs

In my experience, unique Postgres `id` columns
are most often typed with a `uuid` or sequential `bigint`.

This article discusses another ID technique I first saw at
[Chain](https://github.com/chain/chain/blob/main/core/schema.sql),
which was adapted from
[a technique by Instagram](http://instagram-engineering.tumblr.com/post/10853187575/sharding-ids-at-instagram).

Example IDs:

```
user004X70TYG0204
user004X781M00208
user004X7BF50020A
```

These IDs have these properties:

* Readable
* Sortable
* Shardable
* Compact

## Readable

The IDs contain prefixes.
The prefix could be `user` as above or abbreviated like
`org` not `organization` or
`pmt` not `payment`.

The prefix helps make the ID more readable in logs, admin dashboards,
customer support systems, and elsewhere.

The prefix is a per-table ID,
provided as the argument to the `next_id()` function we'll see later:

```sql
CREATE TABLE teams (
  id   text NOT NULL DEFAULT next_id('team'),
  name text NOT NULL,

  UNIQUE (id)
);

CREATE TABLE users (
  id      text NOT NULL DEFAULT next_id('user'),
  email   text NOT NULL,
  team_id text NOT NULL,

  UNIQUE (id)
);
```

The final 8 bytes of the ID are
[Base 32 Crockford encoded](https://www.crockford.com/base32.html).
This further aids readability
and is more compact than `bigint`s.

## Sortable

The IDs are sortable by time.
You can visually see some similarities in the IDs above,
which were provided in ascending order by time.

This is an advantage over `uuid`s.

## Shardable

TODO

## Compact

TODO: storage efficiency

## Postgres implementation

The `next_id()` function we saw earlier:

```sql
CREATE FUNCTION next_id(prefix text)
  RETURNS text
  LANGUAGE plpgsql
  AS $$
  -- Adapted from https://instagram-engineering.com/1cf5a71e5a5c
DECLARE
  our_epoch_ms bigint := 1598390445776; -- do not change
  seq_id bigint;
  now_ms bigint;     -- from unix epoch, not ours
  shard_id int := 1; -- must be different on each shard
  n bigint;
BEGIN
  SELECT
    nextval('id_seq') % 1024 INTO seq_id;
  SELECT
    FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_ms;
  n := (now_ms - our_epoch_ms) << 23;
  n := n | (shard_id << 10);
  n := n | (seq_id);
  RETURN prefix || b32enc_crockford(int8send(n));
END;
$$;
```

`next_id()` depends on an `id_seq` sequence and
`b32enc_crockford()` function:

```sql
CREATE SEQUENCE id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

CREATE FUNCTION b32enc_crockford(src bytea)
  RETURNS text
  LANGUAGE plpgsql
  AS $$
  -- Adapted from the Go package encoding/base32.
  -- See https://golang.org/src/encoding/base32/base32.go
DECLARE
  -- alphabet is base32 alphabet by Douglas Crockford.
  -- It preserves lexical order and avoids visually-similar symbols.
  -- See http://www.crockford.com/wrmg/base32.html
  alphabet text := '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
  dst text := '';
  n integer;
  b0 integer;
  b1 integer;
  b2 integer;
  b3 integer;
  b4 integer;
  b5 integer;
  b6 integer;
  b7 integer;
BEGIN
  FOR r IN 0..(length(src)-1) BY 5
  LOOP
    b0:=0; b1:=0; b2:=0; b3:=0; b4:=0; b5:=0; b6:=0; b7:=0;

    -- Unpack 8x 5-bit source blocks into an 8 byte
    -- destination quantum
    n := length(src) - r;
    IF n >= 5 THEN
      b7 := get_byte(src, r+4) & 31;
      b6 := get_byte(src, r+4) >> 5;
    END IF;
    IF n >= 4 THEN
      b6 := b6 | (get_byte(src, r+3) << 3) & 31;
      b5 := (get_byte(src, r+3) >> 2) & 31;
      b4 := get_byte(src, r+3) >> 7;
    END IF;
    IF n >= 3 THEN
      b4 := b4 | (get_byte(src, r+2) << 1) & 31;
      b3 := (get_byte(src, r+2) >> 4) & 31;
    END IF;
    IF n >= 2 THEN
      b3 := b3 | (get_byte(src, r+1) << 4) & 31;
      b2 := (get_byte(src, r+1) >> 1) & 31;
      b1 := (get_byte(src, r+1) >> 6) & 31;
    END IF;
    b1 := b1 | (get_byte(src, r) << 2) & 31;
    b0 := get_byte(src, r) >> 3;

    -- Encode 5-bit blocks using the base32 alphabet
    dst := dst || substr(alphabet, b0+1, 1);
    dst := dst || substr(alphabet, b1+1, 1);
    IF n >= 2 THEN
      dst := dst || substr(alphabet, b2+1, 1);
      dst := dst || substr(alphabet, b3+1, 1);
    END IF;
    IF n >= 3 THEN
      dst := dst || substr(alphabet, b4+1, 1);
    END IF;
    IF n >= 4 THEN
      dst := dst || substr(alphabet, b5+1, 1);
      dst := dst || substr(alphabet, b6+1, 1);
    END IF;
    IF n >= 5 THEN
      dst := dst || substr(alphabet, b7+1, 1);
    END IF;
  END LOOP;
  RETURN dst;
END;
$$ IMMUTABLE;
```
