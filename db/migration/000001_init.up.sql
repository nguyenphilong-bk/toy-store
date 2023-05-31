CREATE TABLE "users" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR,
    email VARCHAR NOT NULL,
    address VARCHAR,
    birthday DATE,
    phone VARCHAR NOT NULL,
password VARCHAR NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE "orders" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    shipping_fee FLOAT,
    total FLOAT,
    address_shipping VARCHAR,
    user_id uuid NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_orders_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE "brands" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR,
    PRIMARY KEY(id)
);

CREATE TABLE "categories" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR,
    description VARCHAR,
    PRIMARY KEY(id)
);

CREATE TABLE "products" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR,
    origin VARCHAR,
    description VARCHAR,
    image_url VARCHAR,
    price FLOAT,
    stock integer DEFAULT 0,
    PRIMARY KEY(id)
);

CREATE TABLE "product_categories" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    product_id uuid not null,
    category_id uuid not null,
    CONSTRAINT fk_product_category_product FOREIGN KEY(product_id) REFERENCES products(id),
    CONSTRAINT fk_product_category_category FOREIGN KEY(category_id) REFERENCES categories(id)
);

CREATE TABLE "product_brands" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    product_id uuid not null,
    brand_id uuid not null,
    CONSTRAINT fk_product_brand_product FOREIGN KEY(product_id) REFERENCES products(id),
    CONSTRAINT fk_product_brand_brand FOREIGN KEY(brand_id) REFERENCES brands(id)
);

CREATE TABLE "product_orders" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    product_id uuid not null,
    order_id uuid not null,
    CONSTRAINT fk_product_order_product FOREIGN KEY(product_id) REFERENCES products(id),
    CONSTRAINT fk_product_order_order FOREIGN KEY(order_id) REFERENCES orders(id)
);

CREATE TABLE "carts" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    total FLOAT,
    user_id uuid NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_carts_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE "cart_products" (
    id uuid DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    cart_id uuid NOT NULL,
    product_id uuid NOT NULL,
    CONSTRAINT fk_cart_product_cart_id FOREIGN KEY(cart_id) REFERENCES carts(id),
    CONSTRAINT fk_cart_product_product_id FOREIGN KEY(product_id) REFERENCES products(id)
);