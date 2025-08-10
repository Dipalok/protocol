CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  channel TEXT NOT NULL,     -- "email" | "sms"
  sender TEXT,
  recipient TEXT,
  payload JSONB,
  status TEXT NOT NULL DEFAULT 'queued', -- queued, processing, sent, failed
  attempts INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);
