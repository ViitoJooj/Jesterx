import { useEffect, useMemo, useRef, useState } from "react";
import styles from "./RotatingWord.module.scss";

export default function RotatingWord({ items, interval = 2000 }) {
  const safeItems = useMemo(() => items.filter(Boolean), [items]);
  const [current, setCurrent] = useState(0);
  const [previous, setPrevious] = useState(null);
  const timer = useRef(null);
  const currentRef = useRef(0);

  useEffect(() => {
    if (!safeItems.length) return;

    const tick = () => {
      const prev = currentRef.current;
      const next = (prev + 1) % safeItems.length;
      currentRef.current = next;
      setPrevious(prev);
      setCurrent(next);
      timer.current = window.setTimeout(tick, interval);
    };

    timer.current = window.setTimeout(tick, interval);
    return () => clearTimeout(timer.current);
  }, [safeItems, interval]);

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
