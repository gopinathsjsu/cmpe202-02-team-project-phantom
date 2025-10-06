-- Drop existing tables to start fresh (optional, but good for testing)
DROP TABLE IF EXISTS listings;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS LISTING_STATUS;
DROP TYPE IF EXISTS LISTING_CATEGORY;

-- Create custom ENUM types for data integrity
CREATE TYPE LISTING_STATUS AS ENUM ('AVAILABLE', 'SOLD', 'PENDING');
CREATE TYPE LISTING_CATEGORY AS ENUM ('textbooks', 'gadgets', 'essentials', 'other');

-- Create a basic users table (if it doesn't exist from another service)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the listings table
CREATE TABLE listings (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL, -- Price in the smallest currency unit (e.g., cents)
    category LISTING_CATEGORY NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    status LISTING_STATUS DEFAULT 'AVAILABLE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add an index for faster category searches
CREATE INDEX IF NOT EXISTS idx_listings_category ON listings(category);

------------------------------------------------------------------
--                SEED DATA FOR TESTING                         --
------------------------------------------------------------------

-- Insert a sample user
INSERT INTO users (id, email) VALUES (1, 'seller@sjsu.edu')
ON CONFLICT (id) DO NOTHING;

-- Insert sample listings
INSERT INTO listings (id, title, description, price, category, user_id, status) VALUES
(1, 'CMPE202 Advanced Algorithms Textbook', 'Barely used textbook for the CMPE202 class. No highlighting.', 4500, 'textbooks', 1, 'AVAILABLE'),
(2, 'Small Microwave for Dorm Room', 'Works perfectly, great for heating up snacks. Moving out and need to sell.', 2500, 'gadgets', 1, 'AVAILABLE'),
(3, 'Official SJSU Hoodie - Large', 'Official university hoodie, worn a few times. Very comfortable.', 1500, 'essentials', 1, 'AVAILABLE'),
(4, 'Logitech MX Master 3 Mouse', 'Best mouse for productivity. In great condition with original box.', 6000, 'gadgets', 1, 'AVAILABLE'),
(5, 'CHEM1A General Chemistry Textbook', 'Required textbook for the introductory chemistry course. Latest edition.', 3500, 'textbooks', 1, 'SOLD') -- This item is sold and should not appear in search results
ON CONFLICT (id) DO NOTHING;


-- Reset sequences to prevent conflicts if we manually insert with IDs
SELECT setval('users_id_seq', (SELECT MAX(id) from "users"));
SELECT setval('listings_id_seq', (SELECT MAX(id) from "listings"));