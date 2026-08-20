export function isDemo(): boolean {
  return import.meta.env.VITE_DEMO_MODE === "true";
}
