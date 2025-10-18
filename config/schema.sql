CREATE TABLE user_info (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    weight FLOAT(53) NOT NULL,
    height BIGINT NOT NULL,
    date DATE NOT NULL,
    age BIGINT NOT NULL,
    user_id UUID NOT NULL
);

CREATE INDEX idx_user_info_user_id ON user_info(user_id);
CREATE INDEX idx_user_info_user_id_date ON user_info(user_id, date DESC);