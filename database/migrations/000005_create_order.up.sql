CREATE TABLE IF NOT EXISTS orders (
     id SERIAL PRIMARY KEY,
     
     client_id INT NOT NULL,
     status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'received', 'preparing', 'ready', 'delivered', 'canceled')),
     
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     deleted_at TIMESTAMP DEFAULT NULL,

     FOREIGN KEY (client_id) REFERENCES clients(id)
);

CREATE TABLE IF NOT EXISTS order_items (
     id SERIAL PRIMARY KEY,

     order_id INT NOT NULL,
     product_id INT NOT NULL,
     quantity INT NOT NULL,
     price DECIMAL(10, 2) NOT NULL,
     
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     deleted_at TIMESTAMP DEFAULT NULL,

     FOREIGN KEY (order_id) REFERENCES orders(id),
     FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS payments (
     id SERIAL PRIMARY KEY,
     order_id INT NOT NULL,
     status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'failed')),
     method VARCHAR(20) NOT NULL CHECK (method IN ('pix', 'credit_card', 'billet', 'qr_code')),
     amount DECIMAL(10, 2) NOT NULL,
     external_reference VARCHAR(100),
     qr_data TEXT,

     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     deleted_at TIMESTAMP DEFAULT NULL,

     FOREIGN KEY (order_id) REFERENCES orders(id)
);
