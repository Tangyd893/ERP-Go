export function isDemo(): boolean {
  // Never enable mock adapter in production builds.
  if (import.meta.env.PROD) {
    return false;
  }
  return import.meta.env.VITE_DEMO_MODE === "true";
}
