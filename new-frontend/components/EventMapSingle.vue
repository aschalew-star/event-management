<template>
  <div class="relative w-full">
    <!-- Map Canvas Container (Set container layout z-0 mapping layout elements safely) -->
    <div 
      ref="mapContainer" 
      class="w-full rounded-lg overflow-hidden border border-gray-300 dark:border-gray-600 z-0" 
      :style="{ height: height }"
    ></div>
    
    <!-- Custom Control Panels bypassing native default setups styling elements nicely -->
    <div class="absolute bottom-3 right-3 z-[1000] flex flex-col gap-1.5">
      <button @click="zoomIn" type="button" class="p-2 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 rounded shadow hover:bg-gray-50 dark:hover:bg-gray-700 text-xs font-bold">
        +
      </button>
      <button @click="zoomOut" type="button" class="p-2 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 rounded shadow hover:bg-gray-50 dark:hover:bg-gray-700 text-xs font-bold">
        −
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'

const props = defineProps({
  latitude: { type: Number, required: true },
  longitude: { type: Number, required: true },
  height: { type: String, default: '256px' },
  mode: { type: String, default: 'display' }
})

const emit = defineEmits<{
  (e: 'update:location', location: { latitude: number; longitude: number; address: string }): void
}>()

const mapContainer = ref<HTMLDivElement | null>(null)
let L: any = null
let map: any = null
let marker: any = null

onMounted(async () => {
  try {
    // 1. Inject Style requirements dynamically if missing from execution context
    if (!document.getElementById('leaflet-css')) {
      const link = document.createElement('link')
      link.id = 'leaflet-css'
      link.rel = 'stylesheet'
      link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
      document.head.appendChild(link)
    }

    // 2. Load Leaflet dynamic assets client-side only safely
    L = await import('leaflet')

    // 3. FIX: Prevent bundler processing rules from dropping standard asset image pointers
    delete L.Icon.Default.prototype._getIconUrl
    L.Icon.Default.mergeOptions({
      iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
      iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
      shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
    })

    const startLat = props.latitude && !isNaN(props.latitude) ? props.latitude : 9.0319
    const startLng = props.longitude && !isNaN(props.longitude) ? props.longitude : 38.7468

    // 4. Construct Leaflet instance
    map = L.map(mapContainer.value!, {
      zoomControl: false
    }).setView([startLat, startLng], props.mode === 'select' ? 13 : 15)

    // 5. Apply Tile Server links using official OpenStreetMap rules
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '© OpenStreetMap contributors'
    }).addTo(map)

    // 6. Append dynamic custom pins
    marker = L.marker([startLat, startLng], {
      draggable: props.mode === 'select'
    }).addTo(map)

    if (props.mode === 'select') {
      map.on('click', (e: any) => {
        handlePositionChange(e.latlng.lat, e.latlng.lng)
      })
      
      marker.on('dragend', () => {
        const position = marker.getLatLng()
        handlePositionChange(position.lat, position.lng)
      })
    }

    // Forces canvas geometry computation fixes to bypass gray loading block issues
    setTimeout(() => {
      if (map) map.invalidateSize()
    }, 250)

  } catch (error) {
    console.error('Leaflet configuration error:', error)
  }
})

const handlePositionChange = async (lat: number, lng: number) => {
  if (!marker || !map) return
  
  marker.setLatLng([lat, lng])
  map.panTo([lat, lng])
  
  try {
    const address = await reverseGeocode(lat, lng)
    emit('update:location', { latitude: lat, longitude: lng, address })
  } catch (error) {
    emit('update:location', { 
      latitude: lat, 
      longitude: lng, 
      address: `${lat.toFixed(4)}, ${lng.toFixed(4)}` 
    })
  }
}

// Free Tokenless Nominatim Reverse Geocoding API handler
const reverseGeocode = async (lat: number, lng: number): Promise<string> => {
  try {
    const response = await $fetch<any>(
      `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${lat}&lon=${lng}`,
      { headers: { 'User-Agent': 'NuxtEventManagementApp' } }
    )
    return response.display_name || `${lat.toFixed(4)}, ${lng.toFixed(4)}`
  } catch {
    return `${lat.toFixed(4)}, ${lng.toFixed(4)}`
  }
}

// Watch for property value modifications from parent form inputs
watch(() => [props.latitude, props.longitude], ([newLat, newLng]) => {
  if (map && marker && newLat && newLng) {
    const currentCenter = map.getCenter()
    if (Math.abs(currentCenter.lat - newLat) > 0.0001 || Math.abs(currentCenter.lng - newLng) > 0.0001) {
      marker.setLatLng([newLat, newLng])
      map.setView([newLat, newLng])
    }
  }
})

const zoomIn = () => { if (map) map.zoomIn() }
const zoomOut = () => { if (map) map.zoomOut() }

onBeforeUnmount(() => {
  if (map) {
    map.remove()
    map = null
    marker = null
  }
})
</script>