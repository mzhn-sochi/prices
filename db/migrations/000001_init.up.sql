CREATE TABLE items
(
    item_id    UUID     default generateUUIDv4(),
    created_at DATETIME default now(),
    name       String,
    unit       String,
    group      String,
    price      double
) ENGINE = MergeTree()
      PRIMARY KEY (item_id);