use condors;

DROP TABLE IF EXISTS observation;
CREATE TABLE observation (
    `id` INT(11) PRIMARY KEY AUTO_INCREMENT,
    `date` DATE,
    `user` VARCHAR(100),
    `region` VARCHAR(100),
    `nbcondors` INT(11)
);