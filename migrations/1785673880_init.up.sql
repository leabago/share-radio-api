CREATE SCHEMA IF NOT EXISTS stations;

CREATE TABLE IF NOT EXISTS stations.genres (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(50),
    display_order INTEGER DEFAULT 0
    );

CREATE INDEX idx_genres_display_order ON stations.genres (display_order);
CREATE UNIQUE INDEX genres_name ON stations.genres (name);

INSERT INTO stations.genres (name, display_order) VALUES
('pop', 1),
('rock', 2),
('hip hop', 3),
('rap', 4),
('electronic', 5),
('R&B', 6),
('country', 7),
('Latin', 8),
('jazz', 9),
('classical', 10),
('house', 11),
('techno', 12),
('reggae', 13),
('blues', 14),
('soul', 15),
('folk', 16),
('metal', 17),
('funk', 18),
('disco', 19),
('EDM', 20),

('alternative rock', 21),
('classic rock', 22),
('hard rock', 23),
('punk rock', 24),
('indie rock', 25),
('dance-pop', 26),
('synth-pop', 27),
('K-pop', 28),
('Latin pop', 29),
('trance', 30),
('dubstep', 31),
('drum and bass', 32),
('ambient', 33),
('trap', 34),
('gangsta rap', 35),
('drill', 36),
('neo soul', 37),
('smooth jazz', 38),
('bebop', 39),
('electric blues', 40),
('bluegrass', 41),
('Americana', 42),
('salsa', 43),
('bachata', 44),
('reggaeton', 45),
('bossa nova', 46),
('samba', 47),
('tango', 48),
('flamenco', 49),
('afrobeat', 50),

('ska', 51),
('dancehall', 52),
('baroque', 53),
('orchestral', 54),
('opera', 55),
('soundtrack', 56),
('video game music', 57),
('comedy', 58),
('gospel', 59),
('acoustic', 60),
('instrumental', 61),
('a cappella', 62),
('synthwave', 63),
('vaporwave', 64),

('heavy metal', 65),
('thrash metal', 66),
('death metal', 67),
('black metal', 68),
('nu metal', 69),
('progressive rock', 70),
('psychedelic rock', 71),
('garage rock', 72),
('folk rock', 73),
('blues rock', 74),
('southern rock', 75),
('other', 76);

CREATE TABLE IF NOT EXISTS stations.languages (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(50),
    display_order INTEGER DEFAULT 0
    );

CREATE INDEX idx_languages_display_order ON stations.languages (display_order);
CREATE UNIQUE INDEX languages_name ON stations.languages (name);

INSERT INTO stations.languages (name, display_order) VALUES
('English', 1),
('Spanish', 2),
('French', 3),
('German', 4),
('Italian', 5),
('Portuguese', 6),
('Russian', 7),
('Mandarin Chinese', 8),
('Japanese', 9),
('Korean', 10),
('Arabic', 11),
('Hindi', 12),
('Turkish', 13),
('Dutch', 14),
('Polish', 15),
('Ukrainian', 16),
('Romanian', 17),
('Greek', 18),
('Hungarian', 19),
('Czech', 20),
('Swedish', 21),
('Norwegian', 22),
('Danish', 23),
('Finnish', 24),
('Hebrew', 25),
('Persian', 26),
('Indonesian', 27),
('Thai', 28),
('Vietnamese', 29),
('Tagalog', 30),
('Swahili', 31),
('Amharic', 32),
('Catalan', 33),
('Serbian', 34),
('Croatian', 35),
('Slovak', 36),
('Slovenian', 37),
('Bulgarian', 38),
('Estonian', 39),
('Latvian', 40),
('Lithuanian', 41),
('Icelandic', 42),
('Afrikaans', 43),
('Malay', 44),
('Latin', 45),
('Esperanto', 46),
('Multilingual', 47),
('other', 48);

CREATE TABLE IF NOT EXISTS stations.radio_stations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(100) NOT NULL,
    url VARCHAR(500) NOT NULL UNIQUE,
    genre UUID REFERENCES stations.genres(id) NOT NULL,
    language UUID REFERENCES stations.languages(id) NOT NULL,
    icon VARCHAR(500),
    description VARCHAR(500),
    is_active BOOLEAN DEFAULT false,
    is_new BOOLEAN DEFAULT true,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX radio_stations_is_active ON stations.radio_stations (is_active);
CREATE INDEX radio_stations_genre ON stations.radio_stations (genre);
CREATE INDEX radio_stations_language ON stations.radio_stations (language);
CREATE UNIQUE INDEX radio_stations_name ON stations.radio_stations (name);
