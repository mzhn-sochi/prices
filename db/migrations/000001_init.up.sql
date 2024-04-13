CREATE TABLE IF NOT EXISTS items
(
    item_id    UUID     default generateUUIDv4(),
    created_at DATETIME default now(),
    name       String,
    unit       String,
    group      String,
    price      double
) ENGINE = MergeTree()
      PRIMARY KEY (item_id);

CREATE TABLE IF NOT EXISTS shops
(
    district       String,
    businessEntity String,
    inn            String,
    conclusionAt   DATETIME,
    subjectsCount  UInt8,
    addresses      String,
    productNames   String
) ENGINE = ReplacingMergeTree() order by district;