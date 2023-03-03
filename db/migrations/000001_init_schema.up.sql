CREATE TABLE IF NOT EXISTS "song_infos" (
  "id" TEXT PRIMARY KEY,
  "title" TEXT NOT NULL,
  "duration" BIGINT
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'Arctic Monkeys - My Propeller', 305
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'Arctic Monkeys - Crying Lightning', 224
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'Arctic Monkeys - Dance Little Liar', 283
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'Bill Evans - Waltz For Debby', 79
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'Monica Zetterlund, Bill Evans - It Could Happen To You', 179
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'Monica Zetterlund, Bill Evans - Lucky To Be Me', 216
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'Monica Zetterlund, Bill Evans - Come Rain Or Come Shine', 362
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'The Bird And The Bee - My Fair Lady', 214
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'The Bird And The Bee - La La La', 199
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'The Bird And The Bee - Birds And The Bees', 229
);

INSERT INTO song_infos (id, title, duration) (
    select gen_random_uuid(), 'The Bird And The Bee - Again & Again', 165
);
