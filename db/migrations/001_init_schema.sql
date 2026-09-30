-- TRIPO — initial schema (PostgreSQL)
-- Runs automatically on first container init (mounted into /docker-entrypoint-initdb.d).

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    display_name    VARCHAR(100),
    avatar_url      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ITINERARIES (status drives Home vs History screens)
CREATE TYPE itinerary_status AS ENUM ('draft', 'planned', 'ongoing', 'completed', 'cancelled');

CREATE TABLE itineraries (
    id                     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id                UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title                  VARCHAR(150) NOT NULL,
    start_date             DATE NOT NULL,
    end_date               DATE NOT NULL,
    total_budget_estimate  NUMERIC(12,2) DEFAULT 0,
    status                 itinerary_status NOT NULL DEFAULT 'draft',
    cover_image_url        TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_itineraries_user ON itineraries(user_id);
CREATE INDEX idx_itineraries_status ON itineraries(user_id, status);

-- CHAT CONVERSATIONS + MESSAGES
-- Separate conversations per topic; itinerary_id is set only once a trip is generated.
CREATE TABLE chat_conversations (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    itinerary_id    UUID REFERENCES itineraries(id) ON DELETE SET NULL,
    topic           VARCHAR(150) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_chat_conversations_user ON chat_conversations(user_id, updated_at DESC);

CREATE TYPE chat_role AS ENUM ('user', 'assistant', 'system');

CREATE TABLE chat_messages (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES chat_conversations(id) ON DELETE CASCADE,
    role            chat_role NOT NULL,
    content         TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_chat_messages_conversation ON chat_messages(conversation_id, created_at);

-- ITINERARY DAYS
CREATE TABLE itinerary_days (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_id    UUID NOT NULL REFERENCES itineraries(id) ON DELETE CASCADE,
    day_number      INT NOT NULL,
    day_date        DATE NOT NULL,
    UNIQUE (itinerary_id, day_number)
);

-- PLACES (local cache of Google Places data)
CREATE TABLE places (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(200) NOT NULL,
    address         TEXT,
    latitude        DOUBLE PRECISION NOT NULL,
    longitude       DOUBLE PRECISION NOT NULL,
    category        VARCHAR(80),
    google_place_id VARCHAR(150) UNIQUE,
    image_url       TEXT,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_places_coords ON places(latitude, longitude);

-- ITINERARY STOPS (ordered stops per day)
CREATE TYPE stop_status AS ENUM ('pending', 'en_route', 'arrived', 'skipped');

CREATE TABLE itinerary_stops (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_day_id UUID NOT NULL REFERENCES itinerary_days(id) ON DELETE CASCADE,
    place_id         UUID NOT NULL REFERENCES places(id),
    stop_order       INT NOT NULL,
    planned_time     TIME,
    budget_estimate  NUMERIC(10,2) DEFAULT 0,
    status           stop_status NOT NULL DEFAULT 'pending',
    arrived_at       TIMESTAMPTZ,
    departed_at      TIMESTAMPTZ,
    rating           SMALLINT CHECK (rating BETWEEN 1 AND 5),
    actual_price     NUMERIC(10,2),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (itinerary_day_id, stop_order)
);

CREATE INDEX idx_stops_day ON itinerary_stops(itinerary_day_id);

-- ROUTE LEGS (cached pathfinding output; recompute when stops change)
CREATE TABLE route_legs (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_day_id UUID NOT NULL REFERENCES itinerary_days(id) ON DELETE CASCADE,
    from_stop_id     UUID NOT NULL REFERENCES itinerary_stops(id),
    to_stop_id       UUID NOT NULL REFERENCES itinerary_stops(id),
    distance_meters  INT,
    duration_seconds INT,
    polyline         TEXT,
    computed_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_route_legs_day ON route_legs(itinerary_day_id);

-- MEMORIES (photo evidence per stop)
CREATE TABLE memories (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_stop_id UUID NOT NULL REFERENCES itinerary_stops(id) ON DELETE CASCADE,
    photo_url         TEXT NOT NULL,
    caption           VARCHAR(300),
    taken_at          TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_memories_stop ON memories(itinerary_stop_id);

-- ACTIVITIES (event log per stop)
CREATE TABLE activities (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_stop_id UUID NOT NULL REFERENCES itinerary_stops(id) ON DELETE CASCADE,
    description       VARCHAR(300) NOT NULL,
    occurred_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_activities_stop ON activities(itinerary_stop_id);

-- ITINERARY SHARES
CREATE TABLE itinerary_shares (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_id    UUID NOT NULL REFERENCES itineraries(id) ON DELETE CASCADE,
    share_image_url TEXT,
    share_slug      VARCHAR(50) UNIQUE,
    platform        VARCHAR(50),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
