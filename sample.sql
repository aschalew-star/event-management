-- ============================================
-- SAMPLE DATA INSERT FOR ALL TABLES
-- ============================================

-- 1. USERS
INSERT INTO public.users (id, email, password, name, role, avatar_url, bio)
VALUES 
  (gen_random_uuid(), 'john.doe@example.com', 'hashed_password_123', 'John Doe', 'user', 'https://ui-avatars.com/api/?name=John+Doe&background=6366f1&color=fff', 'Event enthusiast and community organizer'),
  (gen_random_uuid(), 'jane.smith@example.com', 'hashed_password_456', 'Jane Smith', 'user', 'https://ui-avatars.com/api/?name=Jane+Smith&background=8b5cf6&color=fff', 'Love attending tech conferences and networking'),
  (gen_random_uuid(), 'mike.johnson@example.com', 'hashed_password_789', 'Mike Johnson', 'event_organizer', 'https://ui-avatars.com/api/?name=Mike+Johnson&background=ec4899&color=fff', 'Professional event organizer with 10+ years experience'),
  (gen_random_uuid(), 'sarah.wilson@example.com', 'hashed_password_101', 'Sarah Wilson', 'user', 'https://ui-avatars.com/api/?name=Sarah+Wilson&background=f59e0b&color=fff', 'Music festival lover and foodie'),
  (gen_random_uuid(), 'alex.chen@example.com', 'hashed_password_202', 'Alex Chen', 'user', 'https://ui-avatars.com/api/?name=Alex+Chen&background=10b981&color=fff', 'Tech entrepreneur and startup enthusiast'),
  (gen_random_uuid(), 'emma.davis@example.com', 'hashed_password_303', 'Emma Davis', 'event_organizer', 'https://ui-avatars.com/api/?name=Emma+Davis&background=3b82f6&color=fff', 'Creative event planner specializing in art exhibitions'),
  (gen_random_uuid(), 'admin@example.com', 'hashed_password_admin', 'Admin User', 'admin', 'https://ui-avatars.com/api/?name=Admin+User&background=ef4444&color=fff', 'Platform administrator');

-- Store user IDs for reference
DO $$
DECLARE
  user1_id UUID := (SELECT id FROM public.users WHERE email = 'john.doe@example.com');
  user2_id UUID := (SELECT id FROM public.users WHERE email = 'jane.smith@example.com');
  user3_id UUID := (SELECT id FROM public.users WHERE email = 'mike.johnson@example.com');
  user4_id UUID := (SELECT id FROM public.users WHERE email = 'sarah.wilson@example.com');
  user5_id UUID := (SELECT id FROM public.users WHERE email = 'alex.chen@example.com');
  user6_id UUID := (SELECT id FROM public.users WHERE email = 'emma.davis@example.com');
  
  -- Category IDs
  cat_music_id UUID;
  cat_tech_id UUID;
  cat_art_id UUID;
  cat_sports_id UUID;
  cat_food_id UUID;
  cat_business_id UUID;
  cat_education_id UUID;
  cat_entertainment_id UUID;
  
  -- Event IDs
  event1_id UUID;
  event2_id UUID;
  event3_id UUID;
  event4_id UUID;
  event5_id UUID;
  event6_id UUID;
  event7_id UUID;
  event8_id UUID;
BEGIN

-- 2. CATEGORIES
INSERT INTO public.categories (id, name, description, icon, color)
VALUES 
  (gen_random_uuid(), 'Music', 'Live music performances, concerts, and music festivals featuring various genres', '🎵', '#6366F1'),
  (gen_random_uuid(), 'Tech', 'Technology conferences, hackathons, coding workshops, and tech meetups', '💻', '#8B5CF6'),
  (gen_random_uuid(), 'Art', 'Art exhibitions, gallery openings, creative workshops, and cultural events', '🎨', '#EC4899'),
  (gen_random_uuid(), 'Sports', 'Sporting events, tournaments, fitness classes, and outdoor activities', '⚽', '#10B981'),
  (gen_random_uuid(), 'Food', 'Food festivals, cooking classes, wine tastings, and culinary experiences', '🍽️', '#F59E0B'),
  (gen_random_uuid(), 'Business', 'Business conferences, networking events, seminars, and professional development', '💼', '#3B82F6'),
  (gen_random_uuid(), 'Education', 'Workshops, lectures, training sessions, and educational programs', '📚', '#EF4444'),
  (gen_random_uuid(), 'Entertainment', 'Comedy shows, theater performances, movie screenings, and entertainment events', '🎭', '#14B8A6');

-- Get category IDs
SELECT id INTO cat_music_id FROM public.categories WHERE name = 'Music';
SELECT id INTO cat_tech_id FROM public.categories WHERE name = 'Tech';
SELECT id INTO cat_art_id FROM public.categories WHERE name = 'Art';
SELECT id INTO cat_sports_id FROM public.categories WHERE name = 'Sports';
SELECT id INTO cat_food_id FROM public.categories WHERE name = 'Food';
SELECT id INTO cat_business_id FROM public.categories WHERE name = 'Business';
SELECT id INTO cat_education_id FROM public.categories WHERE name = 'Education';
SELECT id INTO cat_entertainment_id FROM public.categories WHERE name = 'Entertainment';

-- 3. TAGS
INSERT INTO public.tags (id, name)
VALUES 
  (gen_random_uuid(), 'family-friendly'),
  (gen_random_uuid(), 'outdoor'),
  (gen_random_uuid(), 'indoor'),
  (gen_random_uuid(), 'free'),
  (gen_random_uuid(), 'paid'),
  (gen_random_uuid(), 'workshop'),
  (gen_random_uuid(), 'conference'),
  (gen_random_uuid(), 'festival'),
  (gen_random_uuid(), 'networking'),
  (gen_random_uuid(), 'charity'),
  (gen_random_uuid(), 'virtual'),
  (gen_random_uuid(), 'in-person'),
  (gen_random_uuid(), 'weekend'),
  (gen_random_uuid(), 'weekday'),
  (gen_random_uuid(), 'recurring');

-- 4. EVENTS
-- Event 1: Summer Music Festival (Paid)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'Summer Music Festival 2026',
  'A spectacular 2-day music festival featuring international artists across multiple genres. Enjoy live performances from renowned bands, delicious food trucks, art installations, and a vibrant atmosphere perfect for music lovers of all ages.',
  cat_music_id,
  user3_id,
  49.99,
  false,
  'Central Park Amphitheater',
  '123 Park Avenue, New York, NY 10001',
  40.7829,
  -73.9654,
  '2026-07-15 14:00:00',
  '14:00',
  '23:00',
  'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=800',
  'published',
  245
) RETURNING id INTO event1_id;

-- Event 2: Tech Connect Conference (Paid)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'Tech Connect 2026 - Innovation Summit',
  'Join industry leaders and innovators at the biggest tech conference of the year. Featuring keynotes from CTOs of Fortune 500 companies, hands-on workshops on AI, blockchain, and cloud computing, plus unparalleled networking opportunities with fellow tech enthusiasts.',
  cat_tech_id,
  user5_id,
  299.99,
  false,
  'San Francisco Convention Center',
  '456 Tech Boulevard, San Francisco, CA 94105',
  37.7749,
  -122.4194,
  '2026-08-20 09:00:00',
  '09:00',
  '18:00',
  'https://images.unsplash.com/photo-1540575467063-178a50c2df87?w=800',
  'published',
  189
) RETURNING id INTO event2_id;

-- Event 3: Digital Dreams Art Exhibition (Paid)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'Digital Dreams: Modern Art Exhibition',
  'Experience the intersection of art and technology at this immersive digital art exhibition. Featuring works from 50+ emerging digital artists from around the world, interactive installations, VR experiences, and live digital painting sessions.',
  cat_art_id,
  user6_id,
  15.00,
  false,
  'Downtown Art Gallery',
  '789 Art Street, Austin, TX 78701',
  30.2672,
  -97.7431,
  '2026-06-10 10:00:00',
  '10:00',
  '20:00',
  'https://images.unsplash.com/photo-1533105079780-92b9be482077?w=800',
  'published',
  78
) RETURNING id INTO event3_id;

-- Event 4: Taste of the World (Paid)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'Taste of the World Food Festival',
  'Embark on a culinary journey around the globe without leaving the city. Sample authentic dishes from 30+ countries, watch live cooking demonstrations by award-winning chefs, participate in cooking competitions, and enjoy live entertainment.',
  cat_food_id,
  user3_id,
  35.00,
  false,
  'Waterfront Park',
  '1000 Harbor Drive, Miami, FL 33132',
  25.7617,
  -80.1918,
  '2026-09-05 11:00:00',
  '11:00',
  '21:00',
  'https://images.unsplash.com/photo-1555939594-58d7cb561ad1?w=800',
  'published',
  310
) RETURNING id INTO event4_id;

-- Event 5: City Basketball Championship (Free)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'City Basketball Championship 2026',
  'The most anticipated basketball tournament of the year featuring 16 teams from across the city competing for the championship trophy. Enjoy high-energy games, halftime entertainment, food vendors, and family-friendly activities.',
  cat_sports_id,
  user1_id,
  0,
  true,
  'Sports Arena',
  '555 Stadium Road, Los Angeles, CA 90001',
  34.0522,
  -118.2437,
  '2026-07-25 09:00:00',
  '09:00',
  '17:00',
  'https://images.unsplash.com/photo-1546519638-68e109498ffc?w=800',
  'published',
  167
) RETURNING id INTO event5_id;

-- Event 6: UX Design Workshop (Free)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'UX Design Workshop: From Idea to Prototype',
  'A hands-on workshop designed for beginners and intermediate designers. Learn the fundamentals of user experience design, conduct user research, create wireframes, and build interactive prototypes using industry-standard tools.',
  cat_education_id,
  user6_id,
  0,
  true,
  'Creative Co-working Space',
  '123 Innovation Lane, Chicago, IL 60601',
  41.8781,
  -87.6298,
  '2026-10-15 13:00:00',
  '13:00',
  '17:00',
  'https://images.unsplash.com/photo-1586717791821-3f44a563fa4c?w=800',
  'published',
  92
) RETURNING id INTO event6_id;

-- Event 7: Business Leaders Summit (Paid)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'Business Leaders Summit 2026',
  'An exclusive gathering of business leaders, entrepreneurs, and industry experts. Gain insights into emerging market trends, learn from successful founders, and build valuable business connections in a premium networking environment.',
  cat_business_id,
  user5_id,
  199.99,
  false,
  'Grand Hyatt Hotel',
  '789 Business Avenue, Boston, MA 02110',
  42.3601,
  -71.0589,
  '2026-11-05 08:30:00',
  '08:30',
  '17:30',
  'https://images.unsplash.com/photo-1511578314322-379afb476865?w=800',
  'published',
  134
) RETURNING id INTO event7_id;

-- Event 8: Comedy Night Special (Paid)
INSERT INTO public.events (
  id, title, description, category_id, user_id, price, is_free,
  venue, address, latitude, longitude, event_date, start_time, end_time,
  featured_image, status, view_count
) VALUES (
  gen_random_uuid(),
  'Laugh Factory: Comedy Night Special',
  'An unforgettable evening of laughter with top comedians from around the country. Featuring both established headliners and up-and-coming talent, this comedy show promises a night of non-stop entertainment and good vibes.',
  cat_entertainment_id,
  user3_id,
  25.00,
  false,
  'Comedy Club Downtown',
  '456 Entertainment Street, Las Vegas, NV 89101',
  36.1699,
  -115.1398,
  '2026-12-10 19:00:00',
  '19:00',
  '22:00',
  'https://images.unsplash.com/photo-1527224857830-43a7acc85260?w=800',
  'published',
  56
) RETURNING id INTO event8_id;

-- 5. EVENT IMAGES
-- Add featured images for events
INSERT INTO public.event_images (id, event_id, image_url, is_featured)
SELECT 
  gen_random_uuid(),
  event_id,
  image_url,
  true
FROM (
  VALUES 
    (event1_id, 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=800'),
    (event1_id, 'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=800'),
    (event2_id, 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?w=800'),
    (event2_id, 'https://images.unsplash.com/photo-1505373877841-8d25f7d46678?w=800'),
    (event3_id, 'https://images.unsplash.com/photo-1533105079780-92b9be482077?w=800'),
    (event3_id, 'https://images.unsplash.com/photo-1513364776144-60967b0f800f?w=800'),
    (event4_id, 'https://images.unsplash.com/photo-1555939594-58d7cb561ad1?w=800'),
    (event4_id, 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=800'),
    (event5_id, 'https://images.unsplash.com/photo-1546519638-68e109498ffc?w=800'),
    (event6_id, 'https://images.unsplash.com/photo-1586717791821-3f44a563fa4c?w=800'),
    (event7_id, 'https://images.unsplash.com/photo-1511578314322-379afb476865?w=800'),
    (event8_id, 'https://images.unsplash.com/photo-1527224857830-43a7acc85260?w=800')
) AS images(event_id, image_url)
ON CONFLICT DO NOTHING;

-- Add non-featured images
INSERT INTO public.event_images (id, event_id, image_url, is_featured)
SELECT 
  gen_random_uuid(),
  event_id,
  image_url,
  false
FROM (
  VALUES 
    (event1_id, 'https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=800'),
    (event1_id, 'https://images.unsplash.com/photo-1533174072545-7a4b6ad7a6c3?w=800'),
    (event2_id, 'https://images.unsplash.com/photo-1551817958-5b075e8b9f8c?w=800'),
    (event3_id, 'https://images.unsplash.com/photo-1563299796-17596ed6b017?w=800'),
    (event4_id, 'https://images.unsplash.com/photo-1424847651672-bf20a4b0982b?w=800'),
    (event5_id, 'https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=800')
) AS images(event_id, image_url);

-- 6. EVENT TAGS
INSERT INTO public.event_tags (event_id, tag_id)
SELECT 
  e.id,
  t.id
FROM public.events e
CROSS JOIN public.tags t
WHERE 
  (e.title LIKE '%Music%' AND t.name IN ('outdoor', 'festival', 'paid', 'in-person', 'weekend')) OR
  (e.title LIKE '%Tech%' AND t.name IN ('indoor', 'conference', 'paid', 'networking', 'in-person', 'weekday')) OR
  (e.title LIKE '%Art%' AND t.name IN ('indoor', 'workshop', 'paid', 'in-person', 'weekend')) OR
  (e.title LIKE '%Food%' AND t.name IN ('outdoor', 'festival', 'paid', 'family-friendly', 'in-person', 'weekend')) OR
  (e.title LIKE '%Basketball%' AND t.name IN ('outdoor', 'sports', 'free', 'family-friendly', 'in-person', 'weekend')) OR
  (e.title LIKE '%Workshop%' AND t.name IN ('indoor', 'workshop', 'free', 'in-person', 'weekday')) OR
  (e.title LIKE '%Business%' AND t.name IN ('indoor', 'conference', 'paid', 'networking', 'in-person', 'weekday')) OR
  (e.title LIKE '%Comedy%' AND t.name IN ('indoor', 'paid', 'in-person', 'weekend'));

-- 7. BOOKMARKS (Users bookmarking events they're interested in)
INSERT INTO public.bookmarks (id, user_id, event_id)
SELECT 
  gen_random_uuid(),
  u.id,
  e.id
FROM public.users u
CROSS JOIN public.events e
WHERE 
  (u.email = 'john.doe@example.com' AND e.id IN (event2_id, event4_id)) OR
  (u.email = 'jane.smith@example.com' AND e.id IN (event3_id, event1_id)) OR
  (u.email = 'sarah.wilson@example.com' AND e.id IN (event1_id, event4_id, event8_id)) OR
  (u.email = 'alex.chen@example.com' AND e.id IN (event2_id, event7_id, event6_id)) OR
  (u.email = 'emma.davis@example.com' AND e.id IN (event3_id, event6_id))
ON CONFLICT DO NOTHING;

-- 8. FOLLOWS (Users following events)
INSERT INTO public.follows (id, user_id, event_id)
SELECT 
  gen_random_uuid(),
  u.id,
  e.id
FROM public.users u
CROSS JOIN public.events e
WHERE 
  (u.email = 'john.doe@example.com' AND e.id IN (event1_id, event2_id, event5_id)) OR
  (u.email = 'jane.smith@example.com' AND e.id IN (event1_id, event3_id, event4_id)) OR
  (u.email = 'sarah.wilson@example.com' AND e.id IN (event1_id, event4_id)) OR
  (u.email = 'alex.chen@example.com' AND e.id IN (event2_id, event7_id)) OR
  (u.email = 'emma.davis@example.com' AND e.id IN (event3_id, event6_id))
ON CONFLICT DO NOTHING;

-- 9. TICKETS
-- Confirmed tickets
INSERT INTO public.tickets (id, event_id, user_id, quantity, total_price, status, payment_id, transaction_ref, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  event_id,
  user_id,
  quantity,
  total_price,
  'confirmed',
  payment_id,
  transaction_ref,
  now() - (random() * INTERVAL '30 days'),
  now() - (random() * INTERVAL '30 days')
FROM (
  VALUES 
    (event1_id, user1_id, 2, 99.98, 'pay_123456', 'txn_abcdef'),
    (event1_id, user2_id, 1, 49.99, 'pay_234567', 'txn_bcdefg'),
    (event2_id, user4_id, 1, 299.99, 'pay_345678', 'txn_cdefgh'),
    (event2_id, user6_id, 2, 599.98, 'pay_456789', 'txn_defghi'),
    (event3_id, user2_id, 1, 15.00, 'pay_567890', 'txn_efghij'),
    (event3_id, user5_id, 2, 30.00, 'pay_678901', 'txn_fghijk'),
    (event4_id, user1_id, 3, 105.00, 'pay_789012', 'txn_ghijkl'),
    (event4_id, user5_id, 1, 35.00, 'pay_890123', 'txn_hijklm'),
    (event5_id, user1_id, 4, 0, NULL, NULL),
    (event5_id, user3_id, 2, 0, NULL, NULL),
    (event6_id, user2_id, 1, 0, NULL, NULL),
    (event7_id, user1_id, 1, 199.99, 'pay_901234', 'txn_ijklmn'),
    (event7_id, user4_id, 2, 399.98, 'pay_012345', 'txn_jklmno'),
    (event8_id, user2_id, 2, 50.00, 'pay_123789', 'txn_klmnop')
) AS data(event_id, user_id, quantity, total_price, payment_id, transaction_ref);

-- Pending tickets
INSERT INTO public.tickets (id, event_id, user_id, quantity, total_price, status, payment_id, transaction_ref, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  event_id,
  user_id,
  quantity,
  total_price,
  'pending',
  payment_id,
  transaction_ref,
  now() - (random() * INTERVAL '7 days'),
  now() - (random() * INTERVAL '7 days')
FROM (
  VALUES 
    (event1_id, user5_id, 2, 99.98, 'pay_234890', 'txn_lmnopq'),
    (event2_id, user3_id, 1, 299.99, 'pay_345901', 'txn_mnopqr'),
    (event4_id, user6_id, 1, 35.00, 'pay_456012', 'txn_nopqrs')
) AS data(event_id, user_id, quantity, total_price, payment_id, transaction_ref);

-- Cancelled tickets
INSERT INTO public.tickets (id, event_id, user_id, quantity, total_price, status, payment_id, transaction_ref, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  event_id,
  user_id,
  quantity,
  total_price,
  'cancelled',
  payment_id,
  transaction_ref,
  now() - (random() * INTERVAL '60 days'),
  now() - (random() * INTERVAL '30 days')
FROM (
  VALUES 
    (event1_id, user4_id, 1, 49.99, 'pay_567123', 'txn_opqrst'),
    (event3_id, user1_id, 1, 15.00, 'pay_678234', 'txn_pqrstu')
) AS data(event_id, user_id, quantity, total_price, payment_id, transaction_ref);

-- 10. VERIFY DATA
RAISE NOTICE '✅ Data insertion complete! Here are the counts:';
RAISE NOTICE 'Users: %', (SELECT COUNT(*) FROM public.users);
RAISE NOTICE 'Categories: %', (SELECT COUNT(*) FROM public.categories);
RAISE NOTICE 'Tags: %', (SELECT COUNT(*) FROM public.tags);
RAISE NOTICE 'Events: %', (SELECT COUNT(*) FROM public.events);
RAISE NOTICE 'Event Images: %', (SELECT COUNT(*) FROM public.event_images);
RAISE NOTICE 'Event Tags: %', (SELECT COUNT(*) FROM public.event_tags);
RAISE NOTICE 'Bookmarks: %', (SELECT COUNT(*) FROM public.bookmarks);
RAISE NOTICE 'Follows: %', (SELECT COUNT(*) FROM public.follows);
RAISE NOTICE 'Tickets: %', (SELECT COUNT(*) FROM public.tickets);

END $$;

-- ============================================
-- VIEW SAMPLE DATA
-- ============================================

-- View all events with details
SELECT 
  e.id,
  e.title,
  e.price,
  e.is_free,
  e.event_date,
  e.status,
  c.name as category,
  u.name as organizer,
  e.view_count
FROM public.events e
LEFT JOIN public.categories c ON e.category_id = c.id
LEFT JOIN public.users u ON e.user_id = u.id
ORDER BY e.event_date;

-- View event statistics
SELECT 
  e.title,
  COUNT(DISTINCT b.id) as bookmarks,
  COUNT(DISTINCT f.id) as follows,
  COUNT(DISTINCT t.id) as total_tickets,
  COUNT(DISTINCT CASE WHEN t.status = 'confirmed' THEN t.id END) as confirmed_tickets,
  COALESCE(SUM(CASE WHEN t.status = 'confirmed' THEN t.total_price ELSE 0 END), 0) as revenue
FROM public.events e
LEFT JOIN public.bookmarks b ON e.id = b.event_id
LEFT JOIN public.follows f ON e.id = f.event_id
LEFT JOIN public.tickets t ON e.id = t.event_id
GROUP BY e.id, e.title
ORDER BY revenue DESC;

-- View user activity
SELECT 
  u.name,
  u.email,
  COUNT(DISTINCT b.id) as bookmarks,
  COUNT(DISTINCT f.id) as follows,
  COUNT(DISTINCT t.id) as tickets_purchased,
  COALESCE(SUM(CASE WHEN t.status = 'confirmed' THEN t.total_price ELSE 0 END), 0) as total_spent
FROM public.users u
LEFT JOIN public.bookmarks b ON u.id = b.user_id
LEFT JOIN public.follows f ON u.id = f.user_id
LEFT JOIN public.tickets t ON u.id = t.user_id
GROUP BY u.id, u.name, u.email
ORDER BY total_spent DESC;

-- View ticket sales by category
SELECT 
  c.name as category,
  COUNT(t.id) as tickets_sold,
  COALESCE(SUM(CASE WHEN t.status = 'confirmed' THEN t.total_price ELSE 0 END), 0) as revenue
FROM public.categories c
LEFT JOIN public.events e ON c.id = e.category_id
LEFT JOIN public.tickets t ON e.id = t.event_id
WHERE t.status = 'confirmed'
GROUP BY c.id, c.name
ORDER BY revenue DESC;