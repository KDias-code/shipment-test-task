---
CREATE TABLE IF NOT EXISTS customers (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idn TEXT UNIQUE NOT NULL CHECK (char_length(idn) = 12),
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );