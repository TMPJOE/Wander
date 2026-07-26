import { ref, onMounted } from 'vue'

interface GoogleMapsState {
  isLoaded: boolean
  isLoading: boolean
  error: Error | null
}

const state = ref<GoogleMapsState>({
  isLoaded: false,
  isLoading: false,
  error: null,
})

let loadPromise: Promise<void> | null = null

/**
 * Load the Google Maps JavaScript API script once.
 * Returns loading/error state for components to react to.
 */
export function useGoogleMaps() {
  const apiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY as string | undefined

  onMounted(() => {
    if (!apiKey) {
      state.value.error = new Error('Google Maps API key not configured')
      return
    }

    if (state.value.isLoaded || state.value.isLoading) {
      return
    }

    // Check if already loaded globally
    if (typeof window !== 'undefined' && (window as any).google?.maps) {
      state.value.isLoaded = true
      return
    }

    loadPromise = loadScript(apiKey)
  })

  return {
    isLoaded: () => state.value.isLoaded,
    isLoading: () => state.value.isLoading,
    error: () => state.value.error,
    toMapsUrl,
  }
}

/**
 * Load the Google Maps script dynamically (only once).
 */
async function loadScript(apiKey: string): Promise<void> {
  if (state.value.isLoading || state.value.isLoaded) {
    return loadPromise || Promise.resolve()
  }

  state.value.isLoading = true

  return new Promise((resolve, reject) => {
    // Check again in case of race condition
    if ((window as any).google?.maps) {
      state.value.isLoaded = true
      state.value.isLoading = false
      resolve()
      return
    }

    const script = document.createElement('script')
    script.src = `https://maps.googleapis.com/maps/api/js?key=${encodeURIComponent(apiKey)}&libraries=places`
    script.async = true
    script.defer = true

    script.onload = () => {
      state.value.isLoaded = true
      state.value.isLoading = false
      resolve()
    }

    script.onerror = () => {
      state.value.error = new Error('Failed to load Google Maps script')
      state.value.isLoading = false
      reject(state.value.error)
    }

    document.head.appendChild(script)
  })
}

/**
 * Generate a Google Maps directions URL.
 * @param latitude - Destination latitude
 * @param longitude - Destination longitude
 * @param label - Optional destination label/name
 * @returns URL-safe Google Maps directions link
 */
export function toMapsUrl({
  latitude,
  longitude,
  label,
}: {
  latitude: number | null | undefined
  longitude: number | null | undefined
  label?: string | null
}): string {
  if (latitude != null && longitude != null) {
    const dest = `${latitude},${longitude}`
    const query = label ? encodeURIComponent(label) : dest
    return `https://www.google.com/maps/dir/?api=1&destination=${dest}&destination_place_id=&travelmode=driving&dir_action=navigate`
  }

  // Fallback to text search if no coordinates
  if (label) {
    return `https://www.google.com/maps/search/${encodeURIComponent(label)}`
  }

  return 'https://www.google.com/maps'
}

/**
 * Generate a Google Maps view URL (for display/embedding context).
 */
export function toMapsViewUrl({
  latitude,
  longitude,
  label,
  zoom = 15,
}: {
  latitude: number | null | undefined
  longitude: number | null | undefined
  label?: string | null
  zoom?: number
}): string {
  if (latitude != null && longitude != null) {
    return `https://www.google.com/maps?q=${latitude},${longitude}&z=${zoom}`
  }

  if (label) {
    return `https://www.google.com/maps/search/${encodeURIComponent(label)}`
  }

  return 'https://www.google.com/maps'
}
