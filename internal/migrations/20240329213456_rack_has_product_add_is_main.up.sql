ALTER TABLE rack_has_product
    ADD COLUMN is_main BOOLEAN DEFAULT false;

INSERT INTO rack_has_product(rack_id, product_id, is_main)
VALUES
    (1, 1, true),
    (1, 2, true),
    (2, 3, true),
    (3, 4, true),
    (3, 5, true),
    (3, 6, true);