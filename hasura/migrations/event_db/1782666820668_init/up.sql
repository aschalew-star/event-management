CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SET check_function_bodies = false;
CREATE FUNCTION public.get_event_stats(event_id uuid) RETURNS TABLE(total_bookmarks bigint, total_follows bigint, total_tickets_sold bigint, total_revenue numeric)
    LANGUAGE plpgsql
    AS $_$
BEGIN
    RETURN QUERY
    SELECT 
        (SELECT COUNT(*) FROM bookmarks WHERE event_id = $1) AS total_bookmarks,
        (SELECT COUNT(*) FROM follows WHERE event_id = $1) AS total_follows,
        (SELECT COUNT(*) FROM tickets WHERE event_id = $1 AND status = 'confirmed') AS total_tickets_sold,
        (SELECT COALESCE(SUM(total_price), 0) FROM tickets WHERE event_id = $1 AND status = 'confirmed') AS total_revenue;
END;
$_$;
CREATE FUNCTION public.get_events_near_location(lat numeric, lng numeric, radius_km numeric) RETURNS TABLE(id uuid, title character varying, description text, venue character varying, address text, latitude numeric, longitude numeric, distance numeric)
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
        (6371 * acos(cos(radians(lat)) * cos(radians(e.latitude)) * 
         cos(radians(e.longitude) - radians(lng)) + 
         sin(radians(lat)) * sin(radians(e.latitude)))) AS distance
    FROM events e
    WHERE e.status = 'published'
    AND (6371 * acos(cos(radians(lat)) * cos(radians(e.latitude)) * 
         cos(radians(e.longitude) - radians(lng)) + 
         sin(radians(lat)) * sin(radians(e.latitude)))) < radius_km
    ORDER BY distance;
END;
$$;
CREATE FUNCTION public.increment_event_view_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE events SET view_count = view_count + 1 WHERE id = NEW.event_id;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE TABLE public.bookmarks (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid,
    event_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.categories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    icon character varying(100),
    color character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.event_images (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    event_id uuid,
    image_url character varying(500) NOT NULL,
    is_featured boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.event_tags (
    event_id uuid NOT NULL,
    tag_id uuid NOT NULL
);
CREATE TABLE public.events (
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
    featured_image character varying(500),
    status character varying(50) DEFAULT 'draft'::character varying,
    view_count integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.follows (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid,
    event_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.tags (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.tickets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    event_id uuid,
    user_id uuid,
    quantity integer NOT NULL,
    total_price numeric(10,2) NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying,
    payment_id character varying(255),
    transaction_ref character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.users (
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
ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT bookmarks_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT bookmarks_user_id_event_id_key UNIQUE (user_id, event_id);
ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);
ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.event_images
    ADD CONSTRAINT event_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.event_tags
    ADD CONSTRAINT event_tags_pkey PRIMARY KEY (event_id, tag_id);
ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.follows
    ADD CONSTRAINT follows_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.follows
    ADD CONSTRAINT follows_user_id_event_id_key UNIQUE (user_id, event_id);
ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_name_key UNIQUE (name);
ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
CREATE INDEX idx_bookmarks_event ON public.bookmarks USING btree (event_id);
CREATE INDEX idx_bookmarks_user ON public.bookmarks USING btree (user_id);
CREATE INDEX idx_events_category ON public.events USING btree (category_id);
CREATE INDEX idx_events_date ON public.events USING btree (event_date);
CREATE INDEX idx_events_location ON public.events USING btree (latitude, longitude);
CREATE INDEX idx_events_status ON public.events USING btree (status);
CREATE INDEX idx_events_user ON public.events USING btree (user_id);
CREATE INDEX idx_follows_event ON public.follows USING btree (event_id);
CREATE INDEX idx_follows_user ON public.follows USING btree (user_id);
CREATE INDEX idx_tickets_event ON public.tickets USING btree (event_id);
CREATE INDEX idx_tickets_user ON public.tickets USING btree (user_id);
CREATE TRIGGER increment_event_view AFTER INSERT ON public.tickets FOR EACH ROW EXECUTE FUNCTION public.increment_event_view_count();
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON public.events FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at();
ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT bookmarks_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT bookmarks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.event_images
    ADD CONSTRAINT event_images_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.event_tags
    ADD CONSTRAINT event_tags_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.event_tags
    ADD CONSTRAINT event_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);
ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.follows
    ADD CONSTRAINT follows_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.follows
    ADD CONSTRAINT follows_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
