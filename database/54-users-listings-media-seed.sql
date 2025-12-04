BEGIN;

-- Needed for crypt() / gen_salt()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

WITH
------------------------------------------------------------
-- 1) Build 50 users in the same shape as 51-users-seed.sql
------------------------------------------------------------
new_users AS (
  SELECT
    'user_' || gs               AS user_name,
    'user' || gs || '@sjsu.edu' AS email,
    -- user_1 is admin (role = 1), others normal (role = 0)
    CASE WHEN gs = 1 THEN '0' ELSE '1' END AS role,
    jsonb_build_object(
      'Email', 'user' || gs || '@sjsu.edu'
    ) AS contact
  FROM generate_series(1, 50) AS gs
),

added_users_id AS (
  INSERT INTO users (user_name, email, role, contact)
  SELECT * FROM new_users
  ON CONFLICT (email) DO NOTHING
  RETURNING user_id
),

params AS (
  SELECT crypt('Password123!', gen_salt('bf', 10)) AS pw_hash
),

user_auth_insert AS (
  INSERT INTO user_auth (user_id, password)
  SELECT a.user_id, p.pw_hash
  FROM added_users_id a, params p
  RETURNING user_id
),

------------------------------------------------------------
-- 2) For each inserted user, create exactly 1 listing
------------------------------------------------------------
listing_users AS (
  SELECT
    au.user_id,
    -- Assign a random category per user
    (ARRAY['TEXTBOOK','GADGET','ESSENTIAL','NON-ESSENTIAL','OTHER'])[
      (floor(random() * 5) + 1)::INT
    ] AS category_label
  FROM added_users_id au
),

new_listings AS (
  INSERT INTO listings (title, description, price, category, user_id)
  SELECT
    ------------------------------------------------------------------
    -- Title: depends on category_label
    ------------------------------------------------------------------
    CASE category_label
      WHEN 'TEXTBOOK' THEN format('Textbook: %s',
        (ARRAY[
          'CMPE202 – Software Engineering Materials (Course Pack)',  -- 1 CMPE202-specific item
          'Introduction to Algorithms (CLRS)',
          'Discrete Mathematics and Its Applications',
          'Computer Networking: A Top-Down Approach',
          'Operating Systems: Design and Implementation',
          'Database System Concepts',
          'Computer Organization and Design',
          'Programming Language Principles',
          'Artificial Intelligence: A Modern Approach',
          'Linear Algebra and Its Applications'
        ])[(floor(random() * 10) + 1)::INT]
      )

      WHEN 'GADGET' THEN format('Gadget: %s',
        (ARRAY[
          'Wireless Headphones',
          'Bluetooth Speaker',
          'Portable Power Bank',
          'USB-C Hub for Laptops',
          'Mechanical Keyboard',
          '1080p Webcam for Online Classes'
        ])[(floor(random() * 6) + 1)::INT]
      )

      WHEN 'ESSENTIAL' THEN format('Dorm Essential: %s',
        (ARRAY[
          'Laundry Detergent & Basket',
          'Kitchen Starter Kit (Pots & Pans)',
          'Shower Caddy & Towels',
          'LED Desk Lamp with USB Ports',
          'Set of 30 Hangers',
          'Basic Cleaning Supplies Set'
        ])[(floor(random() * 6) + 1)::INT]
      )

      WHEN 'NON-ESSENTIAL' THEN format('Decor: %s',
        (ARRAY[
          'Wall Poster Pack',
          'Mini Plant Decor Set',
          'LED Strip Lights',
          'Fairy Lights String',
          'Desk Organizer Set',
          'Plush Throw Blanket'
        ])[(floor(random() * 6) + 1)::INT]
      )

      ELSE format('Campus Item: %s',
        (ARRAY[
          'Miscellaneous School Supplies Bundle',
          'Extra Storage Bin',
          'Gym Locker Lock & Accessories',
          'Bike Lock and Light Set',
          'General Student Accessories Pack',
          'Assorted Classroom Essentials'
        ])[(floor(random() * 6) + 1)::INT]
      )
    END AS title,

    ------------------------------------------------------------------
    -- Description: depends on category_label
    ------------------------------------------------------------------
    CASE category_label
      WHEN 'TEXTBOOK' THEN
        'University textbook in good condition. May contain highlights or notes. Suitable for coursework or study reference.'

      WHEN 'GADGET' THEN
        'A reliable electronic gadget in good working condition. Suitable for daily student use, online classes, studying, or entertainment.'

      WHEN 'ESSENTIAL' THEN
        'Everyday dorm essential that helps keep your living space comfortable, organized, or clean.'

      WHEN 'NON-ESSENTIAL' THEN
        'Fun or decorative item to personalize dorm rooms, desks, or shared spaces.'

      ELSE
        'General student item useful for day-to-day campus life or miscellaneous needs.'
    END AS description,

    ------------------------------------------------------------------
    -- Price: random integer between 1000 and 10000 (e.g., cents)
    ------------------------------------------------------------------
    (random() * 9000 + 1000)::INT AS price,

    ------------------------------------------------------------------
    -- Category + owning user
    ------------------------------------------------------------------
    category_label::LISTING_CATEGORY AS category,
    user_id
  FROM listing_users
  RETURNING id, category
)

------------------------------------------------------------
-- 3) Attach 1 stock image per listing based on category
------------------------------------------------------------
INSERT INTO listing_media (listing_id, media_url)
SELECT
  nl.id,
  CASE nl.category
    WHEN 'TEXTBOOK' THEN
      (ARRAY[
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/textbook-1.png',
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/textbook-2.png',
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/textbook-3.png'
      ])[(floor(random() * 3) + 1)::INT]

    WHEN 'GADGET' THEN
      (ARRAY[
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/gadget-1.jpeg',
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/gadget-2.jpeg',
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/gadget-3.jpeg'
      ])[(floor(random() * 3) + 1)::INT]

    WHEN 'ESSENTIAL' THEN
      (ARRAY[
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/essential-1.jpeg',
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/essential-2.jpeg'
      ])[(floor(random() * 2) + 1)::INT]

    WHEN 'NON-ESSENTIAL' THEN
      (ARRAY[
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/nonessential-1.jpeg',
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/nonessential-2.jpeg'
      ])[(floor(random() * 2) + 1)::INT]

    ELSE  -- OTHER
      (ARRAY[
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/other-1.jpeg',
        'https://cmpe202phantomstorage.blob.core.windows.net/phantomstoragecontainer/other-2.jpeg'
      ])[(floor(random() * 2) + 1)::INT]
  END AS media_url
FROM new_listings nl;

COMMIT;
