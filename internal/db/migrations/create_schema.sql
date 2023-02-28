CREATE TABLE lists(
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL
);

CREATE TABLE items(
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_uid UUID REFERENCES lists(uid),
    data VARCHAR NOT NULL
);