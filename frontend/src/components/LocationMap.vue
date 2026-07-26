<template>
  <div class="location-map" :style="{ height: computedHeight }">
    <!-- Google Maps iframe fallback when only label is available -->
    <iframe
      v-if="!hasCoordinates && label"
      :src="iframeSrc"
      width="100%"
      height="100%"
      style="border: 0"
      allowfullscreen
      loading="lazy"
      referrerpolicy="no-referrer-when-downgrade"
      class="location-map__iframe"
    ></iframe>

    <!-- Map container when coordinates are available -->
    <div
      v-else-if="hasCoordinates"
      ref="mapContainer"
      class="location-map__container"
    ></div>

    <!-- Empty state when no location data -->
    <div v-else class="location-map__empty">
      <p>No location information available</p>
    </div>

    <!-- Directions button -->
    <a
      v-if="showDirections && (hasCoordinates || label)"
      :href="directionsUrl"
      target="_blank"
      rel="noopener noreferrer"
      class="location-map__directions"
    >
      <span>Abrir en Google Maps</span>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
        <polyline points="15 3 21 3 21 9" />
        <line x1="10" y1="14" x2="21" y2="3" />
      </svg>
    </a>

    <!-- Fallback notice if maps failed to load -->
    <div v-if="mapsError" class="location-map__error">
      <p>Map unavailable. {{ directionsLabel }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { toMapsUrl, toMapsViewUrl } from '../composables/useGoogleMaps'

const props = withDefaults(
  defineProps<{
    label?: string | null
    latitude?: number | null
    longitude?: number | null
    height?: string
    showDirections?: boolean
  }>(),
  {
    label: null,
    latitude: null,
    longitude: null,
    height: '300px',
    showDirections: true,
  }
)

const mapContainer = ref<HTMLElement | null>(null)
let mapInstance: google.maps.Map | null = null
let markerInstance: google.maps.Marker | null = null
const mapsError = ref<string | null>(null)

const hasCoordinates = computed(() => {
  return props.latitude != null && props.longitude != null
})

const computedHeight = computed(() => props.height)

const iframeSrc = computed(() => {
  if (!props.label) return ''
  const query = encodeURIComponent(props.label)
  return `https://www.google.com/maps/embed/v1/place?key=${encodeURIComponent(
    import.meta.env.VITE_GOOGLE_MAPS_API_KEY || ''
  )}&q=${query}`
})

const directionsUrl = computed(() => {
  return toMapsUrl({
    latitude: props.latitude,
    longitude: props.longitude,
    label: props.label || undefined,
  })
})

const directionsLabel = computed(() => {
  return props.label ? `Search "${props.label}" on Google Maps` : 'Open Google Maps'
})

// Initialize map when coordinates are available and Google Maps is loaded
async function initMap() {
  if (!hasCoordinates.value || !mapContainer.value) return

  // Check if Google Maps is available
  if (typeof window === 'undefined' || !(window as any).google?.maps) {
    // Try to use iframe fallback if we have an API key
    if (import.meta.env.VITE_GOOGLE_MAPS_API_KEY) {
      mapsError.value = 'Interactive map not available'
    }
    return
  }

  try {
    const center = { lat: props.latitude!, lng: props.longitude! }

    mapInstance = new google.maps.Map(mapContainer.value, {
      center,
      zoom: 15,
      disableDefaultUI: false,
      zoomControl: true,
      streetViewControl: false,
      mapTypeControl: false,
      fullscreenControl: true,
    })

    markerInstance = new google.maps.Marker({
      position: center,
      map: mapInstance,
      title: props.label || undefined,
    })
  } catch (e) {
    mapsError.value = 'Failed to initialize map'
    console.error('Map initialization error:', e)
  }
}

// Watch for coordinate changes to update the map
watch(
  () => [props.latitude, props.longitude],
  () => {
    if (mapInstance && hasCoordinates.value) {
      const center = { lat: props.latitude!, lng: props.longitude! }
      mapInstance.setCenter(center)
      if (markerInstance) {
        markerInstance.setPosition(center)
        if (props.label) {
          markerInstance.setTitle(props.label)
        }
      }
    } else if (hasCoordinates.value) {
      // Map wasn't initialized yet, try again
      initMap()
    }
  },
  { immediate: true }
)

onMounted(() => {
  initMap()
})
</script>

<style scoped>
.location-map {
  position: relative;
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
  background-color: #f0f0f0;
}

.location-map__container {
  width: 100%;
  height: 100%;
  min-height: 200px;
}

.location-map__iframe {
  width: 100%;
  height: 100%;
  display: block;
}

.location-map__empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #666;
  font-size: 14px;
}

.location-map__directions {
  position: absolute;
  bottom: 12px;
  right: 12px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background-color: white;
  color: #1a73e8;
  text-decoration: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: background-color 0.2s, transform 0.1s;
  z-index: 10;
}

.location-map__directions:hover {
  background-color: #f8f9fa;
  transform: translateY(-1px);
}

.location-map__directions:active {
  transform: translateY(0);
}

.location-map__error {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.95);
  color: #666;
  font-size: 14px;
  text-align: center;
  padding: 16px;
  z-index: 5;
}

@media (max-width: 640px) {
  .location-map__directions {
    bottom: 8px;
    right: 8px;
    padding: 6px 10px;
    font-size: 12px;
  }
}
</style>
