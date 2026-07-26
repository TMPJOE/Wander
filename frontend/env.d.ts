/// <reference types="vite/client" />
/// <reference types="google.maps" />

// Google Maps TypeScript declarations
declare global {
  interface Window {
    google?: typeof google
  }
}

// Vite environment variables
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_STRIPE_PUBLISHABLE_KEY: string
  readonly VITE_APP_NAME: string
  readonly VITE_APP_ENV: string
  readonly VITE_GOOGLE_MAPS_API_KEY?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
