CREATE TABLE `grade`
(
    `id`         INT(10) UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '自增Id',
    `class`      VARCHAR(128) NOT NULL COMMENT '班级',
    `name`       VARCHAR(128) NOT NULL COMMENT '姓名',
    `score`      INT          NOT NULL COMMENT '分数',
    `subject`    VARCHAR(128) NOT NULL COMMENT '科目',
    `created_at` TIMESTAMP    NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
    `updated_at` TIMESTAMP    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp COMMENT '更新时间'
) COMMENT ='成绩';

CREATE UNIQUE INDEX `unique_idx_grade_class_name_subject` ON `grade` (`class`, `name`, `subject`);
