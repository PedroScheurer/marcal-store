CREATE TABLE tb_product (
                            id          BIGSERIAL    NOT NULL,
                            name        VARCHAR(255) NOT NULL,
                            instructor  VARCHAR(255) NOT NULL,
                            image_url   VARCHAR(255) NOT NULL,
                            video_url   VARCHAR(255) NOT NULL,
                            description VARCHAR(500) NOT NULL,
                            workload    INTEGER      NOT NULL,
                            modules     INTEGER      NOT NULL,
                            price       FLOAT(53)    NOT NULL,
                            currency    VARCHAR(3)   NOT NULL,
                            PRIMARY KEY (id)
);