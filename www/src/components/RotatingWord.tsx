import { useEffect, useMemo, useRef, useState } from "react";
import styles from "../styles/components/RotatingWord.module.scss";

type Props = {
  items: string[];
  interval?: number;
};

export default function RotatingWord({ items, interval = 2000 }: Props) {
  const safeItems = useMemo(() => items.filter(Boolean), [items]);
  const [current, setCurrent] = useState(0);
  const [previous, setPrevious] = useState<number | null>(null);
  const timer = useRef<number | null>(null);

  useEffect(() => {
    if (!safeItems.length) return;
    const tick = () => {
      setPrevious((p) => (p === null ? current : current));
      setCurrent((i) => (i + 1) % safeItems.length);
      timer.current = window.setTimeout(tick, interval);
    };
    timer.current = window.setTimeout(tick, interval);
    return () => {
      if (timer.current) clearTimeout(timer.current);
    };
  }, [safeItems, interval, current]);

  if (!safeItems.length) return null;

  const currentWord = safeItems[current];
  const prevWord = previous !== null ? safeItems[previous] : null;

  return (
    <span className={styles.rotator} aria-live="polite">
      {prevWord !== null && (
        <span key={`out-${previous}`} className={`${styles.word} ${styles.wordOut}`}>
          {prevWord}
        </span>
      )}
      <span key={`in-${current}`} className={`${styles.word} ${styles.wordIn}`}>
        {currentWord}
      </span>
    </span>
  );
}