-- ===================
-- Users
-- ===================
INSERT INTO users (id, name, email, role, active) VALUES 
('1', 'John Doe', 'john@example.com', 'customer', true),
('2', 'Jane Smith', 'jane@example.com', 'admin', true),
('3', 'Bob Wilson', 'bob@example.com', 'customer', true),
('4', 'Alice Johnson', 'alice@example.com', 'customer', true),
('5', 'Mike Chen', 'mike@example.com', 'admin', true),
('6', 'Sarah Wilson', 'sarah@example.com', 'customer', true),
('7', 'David Brown', 'david@example.com', 'customer', false),
('8', 'Emma Davis', 'emma@example.com', 'customer', true),
('9', 'James Miller', 'james@example.com', 'admin', true),
('10', 'Lisa Garcia', 'lisa@example.com', 'customer', true)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Products
-- ===================
INSERT INTO products (id, name, description, price, inventory, available) VALUES
('1', 'MacBook Pro', '14-inch MacBook Pro with M3 chip', 1999.99, 10, true),
('2', 'iPhone 15', 'Latest iPhone with A17 chip', 999.99, 25, true),
('3', 'AirPods Pro', 'Wireless earbuds with noise cancellation', 249.99, 50, true),
('4', 'iPad Pro', '12.9-inch iPad Pro with M2 chip', 1099.99, 15, true),
('5', 'Apple Watch', 'Series 9 GPS + Cellular', 499.99, 30, true),
('6', 'Magic Keyboard', 'Wireless keyboard for Mac', 199.99, 20, true),
('7', 'Studio Display', '27-inch 5K Retina display', 1599.99, 0, false), -- out of stock
('8', 'AirTag', 'Bluetooth tracking device', 29.99, 100, true),
('9', 'HomePod mini', 'Smart speaker with Siri', 99.99, 25, true),
('10', 'Mac Studio', 'Compact pro desktop with M2 Max', 3999.99, 5, true),
('11', 'iPhone 14', 'Previous generation iPhone', 699.99, 40, true),
('12', 'MacBook Air M2', '13-inch lightweight laptop', 1199.99, 12, true),
('13', 'Magic Mouse', 'Wireless multi-touch mouse', 79.99, 35, true),
('14', 'Apple Pencil', '2nd generation stylus for iPad', 129.99, 45, true),
('15', 'Mac mini', 'Compact desktop computer', 599.99, 18, true)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Orders
-- ===================
INSERT INTO orders (id, user_id, quantity, total_price, status, created_at) VALUES
('1', '1', 1, 1999.99, 'completed', NOW()),
('2', '2', 2, 1999.98, 'pending', NOW()),
('3', '3', 1, 249.99, 'shipped', NOW()),
('4', '4', 1, 1099.99, 'shipped', NOW()),
('5', '5', 2, 999.98, 'completed', NOW()),
('6', '6', 1, 199.99, 'pending', NOW()),
('7', '4', 1, 1599.99, 'completed', NOW()),
('8', '8', 4, 119.96, 'shipped', NOW()),
('9', '9', 1, 99.99, 'cancelled', NOW()),
('10', '10', 1, 3999.99, 'pending', NOW())
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Order ↔ Product associations
-- ===================
INSERT INTO order_products (order_id, product_id) VALUES
-- John Doe’s order
('1', '1'),

-- Jane Smith’s order (iPhone 15 + AirPods Pro)
('2', '2'),
('2', '3'),

-- Bob Wilson’s order
('3', '3'),

-- Alice Johnson’s orders
('4', '4'),
('7', '7'),

-- Mike Chen’s order
('5', '5'),

-- Sarah Wilson’s order
('6', '6'),

-- Emma Davis’s order
('8', '8'),

-- James Miller’s order
('9', '9'),

-- Lisa Garcia’s order
('10', '10')
ON CONFLICT DO NOTHING;