-- =========================================================================
-- CORRECTED FUNCTION - get_events_near_location
-- =========================================================================

CREATE OR REPLACE FUNCTION public.get_events_near_location(lat numeric, lng numeric, radius_km numeric) 
RETURNS TABLE(id uuid, title character varying, description text, venue character varying, address text, latitude numeric, longitude numeric, distance numeric)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.title,
        e.description,
        e.venue,
        e.address,
        e.latitude,
        e.longitude,
        CAST((6371 * acos(cos(radians(lat)) * cos(radians(e.latitude)) * cos(radians(e.longitude) - radians(lng)) + 
         sin(radians(lat)) * sin(radians(e.latitude)))) AS numeric) AS distance
    FROM events e
    WHERE e.status = 'published'
    AND (6371 * acos(cos(radians(lat)) * cos(radians(e.latitude)) * cos(radians(e.longitude) - radians(lng)) + 
         sin(radians(lat)) * sin(radians(e.latitude)))) < radius_km
    ORDER BY distance;
END;
$$;

-- =========================================================================
-- SAMPLE DATA FOR ALL TABLES (CORRECTED)
-- =========================================================================

-- 1. USERS
INSERT INTO public.users (id, email, password, name, role, avatar_url, bio) VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'john.doe@example.com', 'hashed_password_123', 'John Doe', 'user', 'https://example.com/avatars/john.jpg', 'Event enthusiast and organizer'),
  ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'jane.smith@example.com', 'hashed_password_456', 'Jane Smith', 'user', 'https://example.com/avatars/jane.jpg', 'Tech conference organizer'),
  ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'mike.johnson@example.com', 'hashed_password_789', 'Mike Johnson', 'admin', 'https://example.com/avatars/mike.jpg', 'Platform administrator'),
  ('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'sarah.wilson@example.com', 'hashed_password_abc', 'Sarah Wilson', 'user', 'https://example.com/avatars/sarah.jpg', 'Music festival organizer'),
  ('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'robert.brown@example.com', 'hashed_password_def', 'Robert Brown', 'user', 'https://example.com/avatars/robert.jpg', 'Sports event planner');

-- 2. CATEGORIES
INSERT INTO public.categories (id, name, description, icon, color) VALUES
  ('f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'Music', 'Music concerts, festivals, and performances', 'music-note', '#FF6B6B'),
  ('a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'Sports', 'Sports events, competitions, and tournaments', 'sports-ball', '#4ECDC4'),
  ('b7eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'Technology', 'Tech conferences, hackathons, and workshops', 'laptop', '#45B7D1'),
  ('c8eebc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'Food', 'Food festivals, cooking classes, and tastings', 'food', '#96CEB4'),
  ('d9eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'Art', 'Art exhibitions, workshops, and galleries', 'palette', '#DDA0DD');

-- 3. EVENTS
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time, 
  status, view_count
) VALUES
  (
    'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Summer Music Festival 2024',
    'A spectacular 3-day music festival featuring top international artists, food trucks, and art installations. Join us for an unforgettable experience!',
    'f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    49.99,
    false,
    'Central Park',
    'Central Park, New York, NY 10001',
    40.7829,
    -73.9654,
    '2024-07-15 10:00:00',
    '10:00',
    '23:00',
    'published',
    1520
  ),
  (
    'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'Tech Conference 2024: Future of AI',
    'Annual technology conference focusing on artificial intelligence, machine learning, and emerging technologies. Keynotes from industry leaders.',
    'b7eebc99-9c0b-4ef8-bb6d-6bb9bd380a88',
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    199.99,
    false,
    'Convention Center',
    '123 Tech Blvd, San Francisco, CA 94105',
    37.7749,
    -122.4194,
    '2024-08-20 09:00:00',
    '09:00',
    '18:00',
    'published',
    2300
  ),
  (
    'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    'Community Yoga in the Park',
    'Free outdoor yoga session for all skill levels. Bring your mat and join the community for a relaxing morning session.',
    'c8eebc99-9c0b-4ef8-bb6d-6bb9bd380a99',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    0.00,
    true,
    'Washington Park',
    'Washington Park, Chicago, IL 60605',
    41.8781,
    -87.6298,
    '2024-06-10 08:00:00',
    '08:00',
    '10:00',
    'published',
    890
  ),
  (
    'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    'International Food Festival',
    'Taste cuisines from around the world. Over 50 food vendors, cooking demonstrations, and live entertainment.',
    'c8eebc99-9c0b-4ef8-bb6d-6bb9bd380a99',
    'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    25.00,
    false,
    'Waterfront Park',
    'Waterfront Park, Miami, FL 33101',
    25.7617,
    -80.1918,
    '2024-09-05 11:00:00',
    '11:00',
    '22:00',
    'published',
    450
  ),
  (
    'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55',
    'Art Exhibition: Modern Masters',
    'Contemporary art exhibition featuring works from emerging and established artists. Opening night with the artists.',
    'd9eebc99-9c0b-4ef8-bb6d-6bb9bd380a00',
    'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55',
    15.00,
    false,
    'Modern Art Gallery',
    '456 Art Avenue, Los Angeles, CA 90001',
    34.0522,
    -118.2437,
    '2024-10-12 18:00:00',
    '18:00',
    '22:00',
    'draft',
    120
  ),
  (
    'd5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66',
    'Marathon 2024',
    'Annual city marathon with full and half marathon options. Fundraising for local charities.',
    'a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77',
    'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55',
    75.00,
    false,
    'City Center',
    'City Center, Boston, MA 02101',
    42.3601,
    -71.0589,
    '2024-11-01 07:00:00',
    '07:00',
    '14:00',
    'published',
    3400
  );

-- 4. EVENT IMAGES
INSERT INTO public.event_images (id, event_id, image_url, is_featured) VALUES
  ('e6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'https://example.com/images/music-festival-banner.jpg', true),
  ('f7eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'https://example.com/images/music-festival-thumb1.jpg', false),
  ('a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'https://example.com/images/music-festival-thumb2.jpg', false),
  ('b9eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'https://example.com/images/tech-conference-banner.jpg', true),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'https://example.com/images/yoga-park-banner.jpg', true),
  ('d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'https://example.com/images/food-festival-banner.jpg', true);

-- 5. TAGS
INSERT INTO public.tags (id, name) VALUES
  ('e2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'rock'),
  ('f3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'jazz'),
  ('a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'outdoor'),
  ('b5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'indoor'),
  ('c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'family-friendly'),
  ('d7eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'workshop'),
  ('e8eebc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'networking'),
  ('f9eebc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'fundraiser');

-- 6. EVENT TAGS
INSERT INTO public.event_tags (event_id, tag_id) VALUES
  ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'e2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33'),
  ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'f3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44'),
  ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55'),
  ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77'),
  ('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'b5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66'),
  ('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'd7eebc99-9c0b-4ef8-bb6d-6bb9bd380a88'),
  ('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'e8eebc99-9c0b-4ef8-bb6d-6bb9bd380a99'),
  ('a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55'),
  ('a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77'),
  ('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55'),
  ('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77'),
  ('d5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55'),
  ('d5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'f9eebc99-9c0b-4ef8-bb6d-6bb9bd380a00');

-- 7. BOOKMARKS
INSERT INTO public.bookmarks (id, user_id, event_id) VALUES
  ('a0febc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22'),
  ('b1febc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44'),
  ('c2febc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
  ('d3febc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'd5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66'),
  ('e4febc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
  ('f5febc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22'),
  ('a6febc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33');

-- 8. FOLLOWS
INSERT INTO public.follows (id, follower_id, followed_user_id) VALUES
  ('b7febc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22'),
  ('c8febc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44'),
  ('d9febc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
  ('e0febc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55'),
  ('f1febc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

-- 9. PAYMENTS
INSERT INTO public.payments (id, user_id, amount, currency, status, transaction_ref, payment_method) VALUES
  ('a2febc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 49.99, 'USD', 'completed', 'TXN-2024-001', 'credit_card'),
  ('b3febc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 25.00, 'USD', 'completed', 'TXN-2024-002', 'paypal'),
  ('c4febc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 199.99, 'USD', 'completed', 'TXN-2024-003', 'credit_card'),
  ('d5febc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 75.00, 'USD', 'pending', 'TXN-2024-004', 'bank_transfer'),
  ('e6febc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 49.99, 'USD', 'completed', 'TXN-2024-005', 'credit_card'),
  ('f7febc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 199.99, 'USD', 'failed', 'TXN-2024-006', 'credit_card'),
  ('a8febc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 15.00, 'USD', 'completed', 'TXN-2024-007', 'paypal');

-- 10. TICKETS
INSERT INTO public.tickets (id, event_id, user_id, payment_id, quantity, total_price, status) VALUES
  ('b9febc99-9c0b-4ef8-bb6d-6bb9bd380a00', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a2febc99-9c0b-4ef8-bb6d-6bb9bd380a33', 2, 99.98, 'confirmed'),
  ('c0febc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b3febc99-9c0b-4ef8-bb6d-6bb9bd380a44', 4, 100.00, 'confirmed'),
  ('d1febc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c4febc99-9c0b-4ef8-bb6d-6bb9bd380a55', 1, 199.99, 'confirmed'),
  ('e2febc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'd5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'd5febc99-9c0b-4ef8-bb6d-6bb9bd380a66', 3, 225.00, 'pending'),
  ('f3febc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'e6febc99-9c0b-4ef8-bb6d-6bb9bd380a77', 1, 49.99, 'confirmed'),
  ('a4febc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'f7febc99-9c0b-4ef8-bb6d-6bb9bd380a88', 2, 399.98, 'cancelled'),
  ('b5febc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'a8febc99-9c0b-4ef8-bb6d-6bb9bd380a99', 2, 30.00, 'confirmed');

-- =========================================================================
-- SAMPLE SELECT QUERIES TO VERIFY DATA
-- =========================================================================

-- Get all events with their categories and organizers
SELECT 
  e.title,
  e.description,
  c.name as category,
  u.name as organizer,
  e.price,
  e.event_date,
  e.status
FROM events e
LEFT JOIN categories c ON e.category_id = c.id
LEFT JOIN users u ON e.user_id = u.id
ORDER BY e.event_date;

-- Get event statistics using the function
SELECT 
  e.title,
  stats.*
FROM events e,
LATERAL get_event_stats(e.id) stats
WHERE e.status = 'published'
LIMIT 5;

-- Find events near a location (New York)
SELECT 
  id,
  title,
  distance,
  venue
FROM get_events_near_location(40.7128, -74.0060, 50)
LIMIT 10;

-- Get all bookmarks for a user with event details
SELECT 
  u.name as user_name,
  e.title as event_title,
  e.event_date
FROM bookmarks b
JOIN users u ON b.user_id = u.id
JOIN events e ON b.event_id = e.id
WHERE u.id = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';

-- Get all followers for a user
SELECT 
  u.name as follower_name,
  u.email as follower_email
FROM follows f
JOIN users u ON f.follower_id = u.id
WHERE f.followed_user_id = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';

-- Get all tickets with payment and event information
SELECT 
  t.id as ticket_id,
  e.title as event_title,
  u.name as attendee_name,
  t.quantity,
  t.total_price,
  t.status as ticket_status,
  p.status as payment_status
FROM tickets t
JOIN events e ON t.event_id = e.id
JOIN users u ON t.user_id = u.id
LEFT JOIN payments p ON t.payment_id = p.id
ORDER BY t.created_at DESC;

-- Get events with their tags
SELECT 
  e.title,
  array_agg(t.name) as tags
FROM events e
JOIN event_tags et ON e.id = et.event_id
JOIN tags t ON et.tag_id = t.id
GROUP BY e.id, e.title;

-- Get total revenue by category
SELECT 
  c.name as category,
  COUNT(t.id) as tickets_sold,
  SUM(t.total_price) as total_revenue
FROM categories c
JOIN events e ON c.id = e.category_id
JOIN tickets t ON e.id = t.event_id
WHERE t.status = 'confirmed'
GROUP BY c.id, c.name
ORDER BY total_revenue DESC;

-- Get most bookmarked events
SELECT 
  e.title,
  COUNT(b.id) as bookmark_count
FROM events e
LEFT JOIN bookmarks b ON e.id = b.event_id
GROUP BY e.id, e.title
ORDER BY bookmark_count DESC
LIMIT 5;