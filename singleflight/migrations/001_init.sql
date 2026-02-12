CREATE TABLE IF NOT EXISTS templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO templates (id, name)
VALUES ('1', 'Template A')
ON CONFLICT (id) DO NOTHING;
