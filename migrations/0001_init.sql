CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    role ENUM('admin', 'buyer') NOT NULL DEFAULT 'buyer',
    status TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    category_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    price DECIMAL(18,2) NOT NULL,
    wholesale_rules JSON NULL,
    stock INT NOT NULL DEFAULT 0,
    status TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_products_category_id (category_id)
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(32) NOT NULL UNIQUE,
    user_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(18,2) NOT NULL,
    total_amount DECIMAL(18,2) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    pay_channel VARCHAR(32) NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    paid_at DATETIME NULL,
    INDEX idx_orders_user_id (user_id),
    INDEX idx_orders_product_id (product_id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_expires_at (expires_at)
);

CREATE TABLE IF NOT EXISTS card_keys (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    order_id BIGINT UNSIGNED NULL,
    locked_at DATETIME NULL,
    sold_at DATETIME NULL,
    INDEX idx_card_keys_product_status_id (product_id, status, id),
    INDEX idx_card_keys_order_id (order_id)
);

CREATE TABLE IF NOT EXISTS payments (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NOT NULL,
    order_no VARCHAR(32) NOT NULL,
    out_trade_no VARCHAR(64) NOT NULL UNIQUE,
    trade_no VARCHAR(64) NULL,
    channel VARCHAR(32) NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    raw_callback JSON NULL,
    paid_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_payments_order_id (order_id),
    INDEX idx_payments_order_no (order_no),
    UNIQUE INDEX uq_payments_channel_trade_no (channel, trade_no)
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NOT NULL,
    card_key_id BIGINT UNSIGNED NOT NULL,
    content_snapshot TEXT NOT NULL,
    delivered_at DATETIME NOT NULL,
    INDEX idx_order_items_order_id (order_id),
    INDEX idx_order_items_card_key_id (card_key_id)
);

CREATE TABLE IF NOT EXISTS product_snapshots (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(128) NOT NULL,
    price DECIMAL(18,2) NOT NULL,
    wholesale_rules JSON NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_product_snapshots_order_id (order_id),
    INDEX idx_product_snapshots_product_id (product_id)
);
