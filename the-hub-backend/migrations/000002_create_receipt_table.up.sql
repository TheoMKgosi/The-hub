CREATE TABLE IF NOT EXISTS receipts (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  image_path TEXT NOT NULL,
  amount DECIMAL(10,2),
  date TIMESTAMP WITH TIME ZONE,
  category_id UUID,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

