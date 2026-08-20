export function vibrate(pattern: number | number[] = 50) {
  try {
    if (typeof navigator !== "undefined" && "vibrate" in navigator) {
      navigator.vibrate(pattern);
    }
  } catch {
    /* unsupported */
  }
}

export function scanFeedback(success: boolean, duplicate = false) {
  if (duplicate) {
    vibrate([30, 50, 30]);
  } else if (success) {
    vibrate(80);
  } else {
    vibrate([80, 100, 80, 100, 200]);
  }
}
