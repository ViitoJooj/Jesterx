-- 0008_StorageSetup.up.sql
-- Creates the Supabase Storage bucket + RLS policies directly in the storage schema.
-- Safe to run multiple times. Skipped silently if the storage schema doesn't exist.

DO $$
BEGIN

-- Only run if this is a Supabase project with the storage extension
IF NOT EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = 'storage' AND table_name = 'buckets'
) THEN
    RAISE NOTICE '0008_StorageSetup: storage schema not found – skipping.';
    RETURN;
END IF;

-- ─── Bucket ───────────────────────────────────────────────────────────────────
INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types, created_at, updated_at)
VALUES (
    'jesterx', 'jesterx',
    true,
    52428800,
    ARRAY[
        'image/jpeg','image/png','image/gif','image/webp','image/svg+xml',
        'video/mp4','video/webm','video/ogg',
        'application/pdf','application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    ],
    NOW(), NOW()
)
ON CONFLICT (id) DO UPDATE
    SET public          = true,
        file_size_limit = EXCLUDED.file_size_limit,
        updated_at      = NOW();

-- ─── RLS policies on storage.objects ─────────────────────────────────────────

-- 1. Authenticated users can upload
DROP POLICY IF EXISTS "jesterx_insert_auth" ON storage.objects;
CREATE POLICY "jesterx_insert_auth" ON storage.objects
    FOR INSERT TO authenticated
    WITH CHECK (bucket_id = 'jesterx');

-- 2. service_role can upload (normally bypasses RLS, but explicit avoids edge cases)
DROP POLICY IF EXISTS "jesterx_insert_service" ON storage.objects;
CREATE POLICY "jesterx_insert_service" ON storage.objects
    FOR INSERT TO service_role
    WITH CHECK (bucket_id = 'jesterx');

-- 3. Public read (anon + authenticated + service_role)
DROP POLICY IF EXISTS "jesterx_select_public" ON storage.objects;
CREATE POLICY "jesterx_select_public" ON storage.objects
    FOR SELECT TO anon, authenticated, service_role
    USING (bucket_id = 'jesterx');

-- 4. Owners and service_role can delete
DROP POLICY IF EXISTS "jesterx_delete_owner" ON storage.objects;
CREATE POLICY "jesterx_delete_owner" ON storage.objects
    FOR DELETE TO authenticated, service_role
    USING (bucket_id = 'jesterx');

RAISE NOTICE '0008_StorageSetup: bucket jesterx + 4 policies created/updated.';

END $$;
