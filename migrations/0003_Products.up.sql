CREATE UNIQUE INDEX idx_products_website_sku  ON products(website_id, sku) WHERE sku IS NOT NULL;
CREATE INDEX         idx_products_website_id   ON products(website_id);
CREATE INDEX         idx_products_website_active ON products(website_id, active);
