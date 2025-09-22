-- =========================
-- Timezones table
-- =========================
CREATE TABLE timezone (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    gmt_offset DECIMAL(4,2) NOT NULL
);

-- =========================
-- Accounts table
-- =========================
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timezone_id INT REFERENCES timezone(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- Identities table
-- =========================
CREATE TABLE identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,                 -- e.g., 'discord' or 'app'
    external_id TEXT NOT NULL,              -- e.g., discord_id or app email
    username TEXT,                          -- snapshot for display purposes
    avatar TEXT,                            -- optional, snapshot of Discord avatar
    password_hash TEXT,                      -- only for app identities
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    UNIQUE(provider, external_id)           -- ensures no duplicate external IDs per provider
);

-- =========================
-- Insert predefined timezones
-- =========================
INSERT INTO timezone (name, gmt_offset) VALUES
('International Date Line West', -12.0),
('Midway Island, Samoa', -11.0),
('Hawaii', -10.0),
('Alaska', -9.0),
('Pacific Time (US & Canada)', -8.0),
('Mountain Time (US & Canada)', -7.0),
('Central Time (US & Canada), Mexico City', -6.0),
('Eastern Time (US & Canada), Bogota, Lima', -5.0),
('Atlantic Time (Canada), Caracas, La Paz', -4.0),
('Newfoundland', -3.5),
('Brazil, Buenos Aires, Georgetown', -3.0),
('Mid-Atlantic', -2.0),
('Azores, Cape Verde Islands', -1.0),
('Western Europe Time, London, Lisbon, Casablanca', 0.0),
('Brussels, Copenhagen, Madrid, Paris', 1.0),
('Kaliningrad, South Africa', 2.0),
('Baghdad, Riyadh, Moscow, St. Petersburg', 3.0),
('Tehran', 3.5),
('Abu Dhabi, Muscat, Baku, Tbilisi', 4.0),
('Kabul', 4.5),
('Ekaterinburg, Islamabad, Karachi, Tashkent', 5.0),
('Bombay, Calcutta, Madras, New Delhi', 5.5),
('Kathmandu', 5.75),
('Almaty, Dhaka, Colombo', 6.0),
('Yangon, Bangkok, Hanoi, Jakarta', 6.5),
('Bangkok, Hanoi, Jakarta', 7.0),
('Beijing, Perth, Singapore, Hong Kong', 8.0),
('Tokyo, Seoul, Osaka, Sapporo, Yakutsk', 9.0),
('Darwin', 9.5),
('Eastern Australia, Guam, Vladivostok', 10.0),
('Magadan, Solomon Islands, New Caledonia', 11.0),
('Auckland, Wellington, Fiji, Kamchatka', 12.0);
