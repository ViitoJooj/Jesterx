-- 0009_FixStorageSetup.up.sql
-- Re-applies correct bucket + RLS policies (0008 had broken current_setting calls).

DO $$
BEGIN

IF NOT EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = 'storage' AND table_name = 'buckets'
) THEN
    RAISE NOTICE '0009_FixStorageSetup: storage schema not present – skipping (use jx storage:setup).';
    RETURN;
END IF;

-- Upsert bucket
INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types, created_at, updated_at)
VALUES (
    'jesterx', 'jesterx', true, 52428800,
    ARRAY[
        'image/jpeg','image/png','image/gif','image/webp','image/svg+xml',
        'video/mp4','video/webm','video/ogg',
        'application/pdf','application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    ],
    NOW(), NOW()
)
ON CONFLICT (id) DO UPDATE
    SET public = true, file_size_limit = EXCLUDED.file_size_limit, updated_at = NOW();

-- Allow authenticated users to upload
DROP POLICY IF EXISTS "jesterx_insert_auth"    ON storage.objects;
DROP POLICY IF EXISTS "jesterx_insert_service" ON storage.objects;
DROP POLICY IF EXISTS "jesterx_select_public"  ON storage.objects;
DROP POLICY IF EXISTS "jesterx_delete_owner"   ON storage.objects;

CREATE POLICY "jesterx_insert_auth" ON storage.objects
    FOR INSERT TO authenticated
    WITH CHECK (bucket_id = 'jesterx');

CREATE POLICY "jesterx_insert_service" ON storage.objects
    FOR INSERT TO service_role
    WITH CHECK (bucket_id = 'jesterx');

CREATE POLICY "jesterx_select_public" ON storage.objects
    FOR SELECT TO anon, authenticated, service_role
    USING (bucket_id = 'jesterx');

CREATE POLICY "jesterx_delete_owner" ON storage.objects
    FOR DELETE TO authenticated, service_role
    USING (bucket_id = 'jesterx');

RAISE NOTICE '0009_FixStorageSetup: bucket jesterx + 4 RLS policies applied.';

END $$;
