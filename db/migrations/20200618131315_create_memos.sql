
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE memos (
    id int AUTO_INCREMENT,
    memo varchar(100),
    PRIMARY KEY(id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE memos;