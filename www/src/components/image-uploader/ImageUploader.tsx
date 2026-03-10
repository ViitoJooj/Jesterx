import { useRef, useState } from "react";
import { API_URL } from "../../hooks/api";
import styles from "./ImageUploader.module.scss";

export type UploadResult = {
  url: string;
  path: string;
  bucket: string;
  mime_type: string;
  size: number;
};

type Props = {
  /** Current image URL to show as preview (controlled) */
  value?: string | null;
  /** Called when a new image is successfully uploaded */
  onUpload: (result: UploadResult) => void;
  /** Optional label for the button */
  label?: string;
  /** Accept attribute for <input type="file"> — defaults to images only */
  accept?: string;
  /** Disable the uploader (e.g. while a parent form is submitting) */
  disabled?: boolean;
  /** Max file size in bytes — defaults to 10MB */
  maxSize?: number;
};

const DEFAULT_ACCEPT = "image/jpeg,image/png,image/gif,image/webp";
const DEFAULT_MAX_SIZE = 10 * 1024 * 1024; // 10 MB

export function ImageUploader({
  value,
  onUpload,
  label = "Enviar imagem",
  accept = DEFAULT_ACCEPT,
  disabled = false,
  maxSize = DEFAULT_MAX_SIZE,
}: Props) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [preview, setPreview] = useState<string | null>(value ?? null);

  // Keep preview in sync if parent changes `value`
  if (value !== undefined && value !== null && value !== preview && !uploading) {
    setPreview(value);
  }

  async function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;

    setError(null);

    // Client-side size validation
    if (file.size > maxSize) {
      setError(`Arquivo muito grande. Máximo: ${Math.round(maxSize / 1024 / 1024)}MB`);
      return;
    }

    // Show local preview immediately while uploading
    const localUrl = URL.createObjectURL(file);
    setPreview(localUrl);

    const form = new FormData();
    form.append("file", file);

    setUploading(true);
    try {
      const res = await fetch(`${API_URL}/api/v1/upload`, {
        method: "POST",
        body: form,
        credentials: "include",
        // Do NOT set Content-Type — browser sets it with the multipart boundary
      });

      const json = await res.json().catch(() => null);

      if (!res.ok || !json?.success) {
        throw new Error(json?.error ?? `Erro ${res.status} ao fazer upload`);
      }

      const result: UploadResult = json.data;
      setPreview(result.url);
      onUpload(result);
    } catch (err: any) {
      setError(err?.message ?? "Falha no upload. Tente novamente.");
      // Revert preview to the previous value
      setPreview(value ?? null);
    } finally {
      setUploading(false);
      // Reset input so same file can be re-selected
      if (inputRef.current) inputRef.current.value = "";
    }
  }

  return (
    <div className={styles.wrapper}>
      {/* Preview */}
      {preview && (
        <div className={styles.previewWrap}>
          <img
            src={preview}
            alt="Preview"
            className={styles.previewImg}
            onError={() => setPreview(null)}
          />
          {uploading && (
            <div className={styles.previewOverlay}>
              <span className={styles.spinner} />
            </div>
          )}
        </div>
      )}

      {/* Hidden file input */}
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        className={styles.hiddenInput}
        onChange={handleChange}
        disabled={disabled || uploading}
        id="image-uploader-input"
      />

      {/* Trigger button */}
      <label
        htmlFor="image-uploader-input"
        className={`${styles.uploadBtn} ${disabled || uploading ? styles.disabled : ""}`}
        aria-disabled={disabled || uploading}
      >
        {uploading ? (
          <>
            <span className={styles.spinnerInline} /> Enviando…
          </>
        ) : (
          <>📷 {label}</>
        )}
      </label>

      {/* Error */}
      {error && <p className={styles.error}>{error}</p>}
    </div>
  );
}
