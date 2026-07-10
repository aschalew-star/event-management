-- Set configuration environment
SET check_function_bodies = false;

-- =========================================================================
-- 1. FUNCTIONS
-- =========================================================================

CREATE OR REPLACE FUNCTION public.get_event_stats(event_id uuid) 
RETURNS TABLE(total_bookmarks bigint, total_tickets_sold bigint, total_revenue numeric)
LANGUAGE plpgsql
AS $_$
BEGIN
    RETURN QUERY
    SELECT 
        (SELECT COUNT(*) FROM bookmarks WHERE bookmarks.event_id = $1) AS total_bookmarks,
        (SELECT COUNT(*) FROM tickets WHERE tickets.event_id = $1 AND tickets.status = 'confirmed') AS total_tickets_sold,
        (SELECT COALESCE(SUM(total_price), 0) FROM tickets WHERE tickets.event_id = $1 AND tickets.status = 'confirmed') AS total_revenue;
END;
$_$;

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
        (6371 * acos(cos(radians(lat)) * cos(radians(e.latitude)) * cos(radians(e.longitude) - radians(lng)) + 
         sin(radians(lat)) * sin(radians(e.latitude)))) AS distance
    FROM events e
    WHERE e.status = 'published'
    AND (6371 * acos(cos(radians(lat)) * cos(radians(e.latitude)) * cos(radians(e.longitude) - radians(lng)) + 
         sin(radians(lat)) * sin(radians(e.latitude)))) < radius_km
    ORDER BY distance;
END;
$$;

CREATE OR REPLACE FUNCTION public.increment_event_view_count() 
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE events SET view_count = view_count + 1 WHERE id = NEW.event_id;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_updated_at() 
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

-- =========================================================================
-- 2. TABLES
-- =========================================================================

CREATE TABLE IF NOT EXISTS public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    role character varying(50) DEFAULT 'user'::character varying,
    avatar_url character varying(500),
    bio text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.categories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    icon character varying(100),
    color character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.events (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(255) NOT NULL,
    description text NOT NULL,
    category_id uuid,
    user_id uuid,
    price numeric(10,2) DEFAULT 0,
    is_free boolean DEFAULT true,
    venue character varying(255),
    address text,
    latitude numeric(10,8),
    longitude numeric(11,8),
    event_date timestamp without time zone NOT NULL,
    start_time character varying(10),
    end_time character varying(10),
    status character varying(50) DEFAULT 'draft'::character varying,
    view_count integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.bookmarks (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    event_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.event_images (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    event_id uuid NOT NULL,
    image_url character varying(500) NOT NULL,
    is_featured boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.tags (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.event_tags (
    event_id uuid NOT NULL,
    tag_id uuid NOT NULL
);

CREATE TABLE IF NOT EXISTS public.follows (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    follower_id uuid NOT NULL,
    followed_user_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.payments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,
    currency character varying(10) DEFAULT 'USD'::character varying,
    status character varying(50) DEFAULT 'pending'::character varying,
    transaction_ref character varying(255),
    payment_method character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.tickets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    event_id uuid NOT NULL,
    user_id uuid NOT NULL,
    payment_id uuid,
    quantity integer NOT NULL,
    total_price numeric(10,2) NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

-- =========================================================================
-- 3. PRIMARY & UNIQUE CONSTRAINTS (Conditional Wrapper)
-- =========================================================================

DO $$
BEGIN
    -- Users
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'users_pkey') THEN
        ALTER TABLE public.users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'users_email_key') THEN
        ALTER TABLE public.users ADD CONSTRAINT users_email_key UNIQUE (email);
    END IF;

    -- Categories
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'categories_pkey') THEN
        ALTER TABLE public.categories ADD CONSTRAINT categories_pkey PRIMARY KEY (id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'categories_name_key') THEN
        ALTER TABLE public.categories ADD CONSTRAINT categories_name_key UNIQUE (name);
    END IF;

    -- Events
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'events_pkey') THEN
        ALTER TABLE public.events ADD CONSTRAINT events_pkey PRIMARY KEY (id);
    END IF;

    -- Bookmarks
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'bookmarks_pkey') THEN
        ALTER TABLE public.bookmarks ADD CONSTRAINT bookmarks_pkey PRIMARY KEY (id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'bookmarks_user_id_event_id_key') THEN
        ALTER TABLE public.bookmarks ADD CONSTRAINT bookmarks_user_id_event_id_key UNIQUE (user_id, event_id);
    END IF;

    -- Event Images
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'event_images_pkey') THEN
        ALTER TABLE public.event_images ADD CONSTRAINT event_images_pkey PRIMARY KEY (id);
    END IF;

    -- Tags
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'tags_pkey') THEN
        ALTER TABLE public.tags ADD CONSTRAINT tags_pkey PRIMARY KEY (id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'tags_name_key') THEN
        ALTER TABLE public.tags ADD CONSTRAINT tags_name_key UNIQUE (name);
    END IF;

    -- Event Tags
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'event_tags_pkey') THEN
        ALTER TABLE public.event_tags ADD CONSTRAINT event_tags_pkey PRIMARY KEY (event_id, tag_id);
    END IF;

    -- Follows
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'follows_pkey') THEN
        ALTER TABLE public.follows ADD CONSTRAINT follows_pkey PRIMARY KEY (id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'follows_follower_followed_key') THEN
        ALTER TABLE public.follows ADD CONSTRAINT follows_follower_followed_key UNIQUE (follower_id, followed_user_id);
    END IF;

    -- Payments
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'payments_pkey') THEN
        ALTER TABLE public.payments ADD CONSTRAINT payments_pkey PRIMARY KEY (id);
    END IF;

    -- Tickets
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'tickets_pkey') THEN
        ALTER TABLE public.tickets ADD CONSTRAINT tickets_pkey PRIMARY KEY (id);
    END IF;
END $$;

-- =========================================================================
-- 4. INDEXES
-- =========================================================================

CREATE INDEX IF NOT EXISTS idx_bookmarks_event ON public.bookmarks USING btree (event_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_user ON public.bookmarks USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_events_category ON public.events USING btree (category_id);
CREATE INDEX IF NOT EXISTS idx_events_date ON public.events USING btree (event_date);
CREATE INDEX IF NOT EXISTS idx_events_location ON public.events USING btree (latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_events_status ON public.events USING btree (status);
CREATE INDEX IF NOT EXISTS idx_events_user ON public.events USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_follows_followed ON public.follows USING btree (followed_user_id);
CREATE INDEX IF NOT EXISTS idx_follows_follower ON public.follows USING btree (follower_id);
CREATE INDEX IF NOT EXISTS idx_payments_user ON public.payments USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_event ON public.tickets USING btree (event_id);
CREATE INDEX IF NOT EXISTS idx_tickets_user ON public.tickets USING btree (user_id);

-- =========================================================================
-- 5. TRIGGERS (Drop if exists to avoid duplication)
-- =========================================================================

DROP TRIGGER IF EXISTS increment_event_view ON public.tickets;
CREATE TRIGGER increment_event_view AFTER INSERT ON public.tickets FOR EACH ROW EXECUTE FUNCTION public.increment_event_view_count();

DROP TRIGGER IF EXISTS update_categories_updated_at ON public.categories;
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();

DROP TRIGGER IF EXISTS update_events_updated_at ON public.events;
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON public.events FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();

DROP TRIGGER IF EXISTS update_payments_updated_at ON public.payments;
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON public.payments FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();

DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();

-- =========================================================================
-- 6. FOREIGN KEY CONSTRAINTS (Conditional Wrapper)
-- =========================================================================

DO $$
BEGIN
    -- Events FKs
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'events_category_id_fkey') THEN
        ALTER TABLE public.events ADD CONSTRAINT events_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'events_user_id_fkey') THEN
        ALTER TABLE public.events ADD CONSTRAINT events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
    END IF;

    -- Bookmarks FKs
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'bookmarks_event_id_fkey') THEN
        ALTER TABLE public.bookmarks ADD CONSTRAINT bookmarks_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'bookmarks_user_id_fkey') THEN
        ALTER TABLE public.bookmarks ADD CONSTRAINT bookmarks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
    END IF;

    -- Event Images FKs
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'event_images_event_id_fkey') THEN
        ALTER TABLE public.event_images ADD CONSTRAINT event_images_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
    END IF;

    -- Event Tags FKs
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'event_tags_event_id_fkey') THEN
        ALTER TABLE public.event_tags ADD CONSTRAINT event_tags_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'event_tags_tag_id_fkey') THEN
        ALTER TABLE public.event_tags ADD CONSTRAINT event_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE;
    END IF;

    -- Follows FKs
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'follows_followed_user_id_fkey') THEN
        ALTER TABLE public.follows ADD CONSTRAINT follows_followed_user_id_fkey FOREIGN KEY (followed_user_id) REFERENCES public.users(id) ON DELETE CASCADE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'follows_follower_id_fkey') THEN
        ALTER TABLE public.follows ADD CONSTRAINT follows_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES public.users(id) ON DELETE CASCADE;
    END IF;

    -- Payments FKs
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'payments_user_id_fkey') THEN
        ALTER TABLE public.payments ADD CONSTRAINT payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;
    END IF;

    -- Tickets FKs
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'tickets_event_id_fkey') THEN
        ALTER TABLE public.tickets ADD CONSTRAINT tickets_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'tickets_payment_id_fkey') THEN
        ALTER TABLE public.tickets ADD CONSTRAINT tickets_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payments(id) ON DELETE SET NULL;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'tickets_user_id_fkey') THEN
        ALTER TABLE public.tickets ADD CONSTRAINT tickets_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
    END IF;
END $$;