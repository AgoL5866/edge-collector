create table block_header_infos
(
    req_time    DateTime,
    duration    Int64,
    region      String,
    chain       String,
    node_alias  String,
    number      Int64,
    timestamp   DateTime,
    hash        String,
    parent_hash String,
    reason      String
)
    ENGINE = MergeTree() ORDER BY (req_time, chain, node_alias, number)
        PARTITION BY (toYYYYMM(req_time), chain, region);