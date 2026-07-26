<template>
  <div class="location-picker-map">
    <!-- Search input for Places Autocomplete -->
    <div v-if="isMapsAvailable" class="location-picker-map__search">
      <label :for="searchInputId" class="location-picker-map__label">
        Buscar lugar
      </label>
      <div class="location-picker-map__input-wrapper">
        <input
          :id="searchInputId"
          ref="searchInput"
          type="text"
          class="location-picker-map__input"
          placeholder="Ingrese dirección o lugar..."
          :value="searchQuery"
          @input="onSearchInput"
        />
        <button
          v-if="searchQuery"
          type="button"
          class="location-picker-map__clear"
          @click="clearLocation"
          aria-label="Limpiar ubicación seleccionada"
        >
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
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Notice when Maps is unavailable -->
    <div v-else class="location-picker-map__notice">
      <p>
        Google Maps no está disponible. Puede ingresar la ubicación manualmente en los campos de texto.
      </p>
    </div>

    <!-- Map container -->
    <div
      ref="mapContainer"
      class="location-picker-map__container"
      :class="{ 'location-picker-map__container--disabled': !isMapsAvailable }"
      @click="onMapClick"
    >
      <div v-if="!hasLocation && isMapsAvailable" class="location-picker-map__placeholder">
        <p>Haga clic en el mapa para seleccionar una ubicación</p>
        <p class="location-picker-map__placeholder--small">
          O use el buscador arriba
        </p>
      </div>
    </div>

    <!-- Selected location info -->
    <div v-if="hasLocation" class="location-picker-map__info">
      <div class="location-picker-map__coords">
        <span class="location-picker-map__coords-label">Coordenadas:</span>
        <span class="location-picker-map__coords-value">
          {{ formattedLatitude }}, {{ formattedLongitude }}
        </span>
      </div>
      <button
        v-if="isMapsAvailable"
        type="button"
        class="location-picker-map__clear-btn"
        @click="clearLocation"
      >
        Limpiar ubicación
      </button>
    </div>

    <!-- Manual entry fields (always available as fallback) -->
    <div class="location-picker-map__manual">
      <p class="location-picker-map__manual-title">Entrada manual</p>
      <div class="location-picker-map__manual-fields">
        <slot name="manual-fields">
          <!-- Parent can provide custom manual fields or use default -->
          <p class="location-picker-map__manual-hint">
            Los campos de ubicación y punto de encuentro pueden editarse directamente en el formulario.
          </p>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    label?: string | null
    latitude?: number | null
    longitude?: number | null
  }>(),
  {
    label: null,
    latitude: null,
    longitude: null,
  }
)

const emit = defineEmits<{
  (e: 'update:latitude', value: number | null): void
  (e: 'update:longitude', value: number | null): void
  (e: 'update:label', value: string | null): void
}>()

const mapContainer = ref<HTMLElement | null>(null)
const searchInput = ref<HTMLInputElement | null>(null)
const searchQuery = ref('')

let mapInstance: google.maps.Map | null = null
let markerInstance: google.maps.Marker | null = null
let autocompleteInstance: google.maps.places.Autocomplete | null = null
let isDragging = false

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY as string | undefined
const isMapsAvailable = computed(() => {
  return !!apiKey && typeof window !== 'undefined' && !!(window as any).google?.maps
})

const hasLocation = computed(() => {
  return props.latitude != null && props.longitude != null
})

const formattedLatitude = computed(() => {
  return props.latitude?.toFixed(6) || ''
})

const formattedLongitude = computed(() => {
  return props.longitude?.toFixed(6) || ''
})

const searchInputId = `location-search-${Math.random().toString(36).slice(2, 9)}`

// Initialize map and autocomplete
async function initMap() {
  if (!isMapsAvailable.value || !mapContainer.value) return

  try {
    const center = hasLocation.value
      ? { lat: props.latitude!, lng: props.longitude! }
      : { lat: -34.6037, lng: -58.3816 } // Default to Buenos Aires

    mapInstance = new google.maps.Map(mapContainer.value, {
      center,
      zoom: hasLocation.value ? 15 : 12,
      disableDefaultUI: false,
      zoomControl: true,
      streetViewControl: false,
      mapTypeControl: false,
      fullscreenControl: true,
    })

    // Create draggable marker
    markerInstance = new google.maps.Marker({
      position: center,
      map: mapInstance,
      draggable: true,
      title: props.label || undefined,
      visible: hasLocation.value,
    })

    // Handle marker drag events
    markerInstance.addListener('dragstart', () => {
      isDragging = true
    })

    markerInstance.addListener('dragend', (event: google.maps.MapMouseEvent) => {
      isDragging = false
      if (event.latLng) {
        updateLocation(event.latLng.lat(), event.latLng.lng())
      }
    })

    // Setup Places Autocomplete
    if (searchInput.value) {
      autocompleteInstance = new google.maps.places.Autocomplete(searchInput.value, {
        fields: ['geometry', 'name', 'formatted_address'],
        types: ['geocode', 'establishment'],
      })

      autocompleteInstance.addListener('place_changed', () => {
        const place = autocompleteInstance?.getPlace()
        if (place?.geometry?.location) {
          const lat = place.geometry.location.lat()
          const lng = place.geometry.location.lng()
          const name = place.name || place.formatted_address || ''
          updateLocation(lat, lng, name)
        }
      })
    }

    // Update marker position if coordinates change externally
    watch(
      () => [props.latitude, props.longitude, props.label],
      () => {
        if (markerInstance && hasLocation.value) {
          const newPosition = { lat: props.latitude!, lng: props.longitude! }
          markerInstance.setPosition(newPosition)
          markerInstance.setVisible(true)
          markerInstance.setTitle(props.label || '')
          mapInstance?.setCenter(newPosition)
          mapInstance?.setZoom(15)
        }
      },
      { immediate: true }
    )
  } catch (e) {
    console.error('Failed to initialize map:', e)
  }
}

// Handle map click to select location
function onMapClick(event: MouseEvent) {
  if (!isMapsAvailable.value || !mapInstance || isDragging) return

  // Only process direct map clicks, not marker drags
  const rect = mapContainer.value?.getBoundingClientRect()
  if (!rect) return

  // Get click position relative to map
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top

  // Use projection to convert pixel to lat/lng
  const projection = mapInstance.getProjection()
  if (!projection) return

  const point = projection.fromContainerPixelToLatLng({ x, y })
  if (point) {
    updateLocation(point.lat(), point.lng())
  }
}

// Update location and emit changes
function updateLocation(lat: number, lng: number, label?: string) {
  emit('update:latitude', lat)
  emit('update:longitude', lng)
  if (label) {
    emit('update:label', label)
  }
}

// Clear selected location
function clearLocation() {
  emit('update:latitude', null)
  emit('update:longitude', null)
  searchQuery.value = ''
  if (markerInstance) {
    markerInstance.setVisible(false)
  }
}

// Handle search input
function onSearchInput(event: Event) {
  const target = event.target as HTMLInputElement
  searchQuery.value = target.value
}

// Sync label to search query when it changes externally
watch(
  () => props.label,
  (newLabel) => {
    if (newLabel && !isDragging) {
      searchQuery.value = newLabel
    }
  },
  { immediate: true }
)

onMounted(() => {
  initMap()
})
</script>

<style scoped>
.location-picker-map {
  width: 100%;
}

.location-picker-map__search {
  margin-bottom: 12px;
}

.location-picker-map__label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.location-picker-map__input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.location-picker-map__input {
  width: 100%;
  padding: 10px 36px 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.location-picker-map__input:focus {
  outline: none;
  border-color: #1a73e8;
  box-shadow: 0 0 0 3px rgba(26, 115, 232, 0.1);
}

.location-picker-map__clear {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: transparent;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  color: #6b7280;
  transition: background-color 0.2s, color 0.2s;
}

.location-picker-map__clear:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.location-picker-map__notice {
  padding: 12px;
  background-color: #fef3c7;
  border: 1px solid #fcd34d;
  border-radius: 6px;
  color: #92400e;
  font-size: 13px;
  margin-bottom: 12px;
}

.location-picker-map__container {
  width: 100%;
  height: 300px;
  border-radius: 8px;
  overflow: hidden;
  background-color: #e5e7eb;
  cursor: pointer;
  position: relative;
}

.location-picker-map__container--disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.location-picker-map__placeholder {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: #6b7280;
  font-size: 14px;
  pointer-events: none;
}

.location-picker-map__placeholder--small {
  font-size: 12px;
  margin-top: 4px;
  color: #9ca3af;
}

.location-picker-map__info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  margin-top: 12px;
}

.location-picker-map__coords {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.location-picker-map__coords-label {
  color: #6b7280;
}

.location-picker-map__coords-value {
  font-family: monospace;
  color: #111827;
}

.location-picker-map__clear-btn {
  padding: 6px 12px;
  background-color: white;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 13px;
  color: #374151;
  cursor: pointer;
  transition: background-color 0.2s, border-color 0.2s;
}

.location-picker-map__clear-btn:hover {
  background-color: #f9fafb;
  border-color: #9ca3af;
}

.location-picker-map__manual {
  margin-top: 16px;
  padding: 12px;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
}

.location-picker-map__manual-title {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.location-picker-map__manual-fields {
  font-size: 13px;
  color: #6b7280;
}

.location-picker-map__manual-hint {
  line-height: 1.5;
}

@media (max-width: 640px) {
  .location-picker-map__container {
    height: 250px;
  }

  .location-picker-map__info {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
