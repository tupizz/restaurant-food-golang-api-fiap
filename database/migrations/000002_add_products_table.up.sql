CREATE TABLE IF NOT EXISTS categories (
     id SERIAL PRIMARY KEY,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     deleted_at TIMESTAMP WITH TIME ZONE,
     name VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO categories (name) VALUES ('Lanche');
INSERT INTO categories (name) VALUES ('Bebida');
INSERT INTO categories (name) VALUES ('Sobremesa');
INSERT INTO categories (name) VALUES ('Acompanhamento');

CREATE TABLE IF NOT EXISTS products (
     id SERIAL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     description VARCHAR(255) NOT NULL,
     price DECIMAL(10, 2) NOT NULL,
     category_id INT NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     deleted_at TIMESTAMP WITH TIME ZONE,
     FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE IF NOT EXISTS products_images (
     id SERIAL PRIMARY KEY,
     product_id INT NOT NULL,
     image VARCHAR(255) NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     deleted_at TIMESTAMP WITH TIME ZONE,
     FOREIGN KEY (product_id) REFERENCES products (id)
);
