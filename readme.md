CREATE TABLE simple_struct (
id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
bool_column BOOLEAN NOT NULL,
bool_ptr BOOLEAN NULL,
int_column INT NOT NULL,
int_ptr INT NULL,

                               int8_column TINYINT NOT NULL,
                               int8_ptr TINYINT NULL,

                               int16_column SMALLINT NOT NULL,
                               int16_ptr SMALLINT NULL,

                               int32_column INT NOT NULL,
                               int32_ptr INT NULL,

                               int64_column BIGINT NOT NULL,
                               int64_ptr BIGINT NULL,

                               uint_column BIGINT UNSIGNED NOT NULL,
                               uint_ptr BIGINT UNSIGNED NULL,

                               uint8_column TINYINT UNSIGNED NOT NULL,
                               uint8_ptr TINYINT UNSIGNED NULL,

                               uint16_column SMALLINT UNSIGNED NOT NULL,
                               uint16_ptr SMALLINT UNSIGNED NULL,

                               uint32_column INT UNSIGNED NOT NULL,
                               uint32_ptr INT UNSIGNED NULL, -- 注意这里应该是Uint32Ptr而不是Uint32Ptr

                               uint64_column BIGINT UNSIGNED NOT NULL,
                               uint64_ptr BIGINT UNSIGNED NULL, -- 同样，这里应该是Uint64Ptr

                               float32_column FLOAT NOT NULL,
                               float32_ptr FLOAT NULL,

                               float64_column DOUBLE NOT NULL,
                               float64_ptr DOUBLE NULL,

                               byte_column TINYINT UNSIGNED NOT NULL,
                               byte_ptr TINYINT UNSIGNED NULL,

                               byte_array TEXT, -- 切片类型通常存储为TEXT或BLOB

                               string_column VARCHAR(255) NOT NULL, -- 根据需要调整长度

                               null_string_ptr VARCHAR(255) NULL,
                               null_int16_ptr SMALLINT NULL,
                               null_int32_ptr INT NULL,
                               null_int64_ptr BIGINT NULL,
                               null_bool_ptr BOOLEAN NULL,
                               null_time_ptr DATETIME NULL,
                               null_float64_ptr DOUBLE NULL,

                               json_column JSON
);

go test ./... 跑单元测试
go test -tags=e2e ./... 单元测试和集成测试一起跑

// go:build e2e