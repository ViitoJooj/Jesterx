import { createClient } from "@supabase/supabase-js";

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL as string;
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY as string;

export const supabase = createClient(supabaseUrl, supabaseAnonKey);

export async function uploadImage(file: File, path: string): Promise<string> {
  const { data, error } = await supabase.storage
    .from("img")
    .upload(path, file, { upsert: true, cacheControl: "3600" });

  if (error) throw new Error(error.message);

  const { data: urlData } = supabase.storage
    .from("img")
    .getPublicUrl(data.path);

  return urlData.publicUrl;
}
