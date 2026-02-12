CREATE TABLE IF NOT EXISTS templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO templates (id, name) VALUES
    ('1', 'Template A'),
    ('2', 'Template B'),
    ('3', 'Template C'),
    ('4', 'Template D'),
    ('5', 'Template E'),
    ('6', 'Template F'),
    ('7', 'Template G'),
    ('8', 'Template H'),
    ('9', 'Template I'),
    ('10', 'Template J')
ON CONFLICT (id) DO NOTHING;
