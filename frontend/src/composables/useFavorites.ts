import { ref, computed } from 'vue'
import { useApi } from './useApi'
import { useAuthState } from './useAuthState'

// Module-scoped singleton state. This mirrors the pattern in useAuthState /
// useApi: a single source of truth for favorite tour IDs shared across every
// view that renders a heart (Explore grid, Favorites tab, TourCard). Module
// scope guarantees one shared set per page so toggling a heart anywhere keeps
// every other card in sync without a reload.
const favoritedIds = ref<Set<number>>(new Set())
const loaded = ref(false)
let loadingPromise: Promise<void> | null = null

async function loadOnce(): Promise<void> {
  const authState = useAuthState()
  // Favorites are user-specific. Skip entirely for anonymous visitors so the
  // set stays empty and hearts render as the un-liked default.
  if (!authState.token.value) {
    favoritedIds.value = new Set()
    loaded.value = true
    return
  }

  if (loadingPromise) return loadingPromise
  const api = useApi()
  loadingPromise = (async () => {
    try {
      const res = await api.get('/favorites')
      const rows = Array.isArray(res.data) ? (res.data as Array<{ id: number }>) : []
      favoritedIds.value = new Set(rows.map((t) => t.id))
    } catch (e) {
      console.error('Failed to load favorites', e)
      favoritedIds.value = new Set()
    } finally {
      loaded.value = true
      loadingPromise = null
    }
  })()
  return loadingPromise
}

export function useFavorites() {
  const api = useApi()
  const authState = useAuthState()

  // hydrate() is idempotent: it loads /favorites once, then no-ops on subsequent
  // calls. Call it from any view that needs the heart state before first paint.
  // Pass { force: true } to re-fetch after out-of-band mutations (e.g. login).
  async function hydrate(opts?: { force?: boolean }): Promise<void> {
    if (opts?.force) {
      loaded.value = false
      favoritedIds.value = new Set()
    }
    if (!loaded.value) {
      await loadOnce()
    }
  }

  function isFavorited(id: number): boolean {
    return favoritedIds.value.has(id)
  }

  // isFavoritedReactive returns a computed so templates stay reactive when the
  // underlying set is replaced (we assign a new Set on each mutation to keep
  // Vue's reactivity simple and avoid Set-vs-reactivity caveats).
  function isFavoritedReactive(id: number) {
    return computed(() => favoritedIds.value.has(id))
  }

  function setFavorited(id: number, value: boolean) {
    const next = new Set(favoritedIds.value)
    if (value) next.add(id)
    else next.delete(id)
    favoritedIds.value = next
  }

  // toggleFavorite performs an optimistic update: flip the set immediately so
  // the UI feels instant, then call the API. On failure it rolls back. Returns
  // the new favorited state (or null on auth failure, meaning "did nothing").
  async function toggleFavorite(id: number): Promise<boolean | null> {
    if (!authState.token.value) return null

    const previous = isFavorited(id)
    setFavorited(id, !previous)
    try {
      if (previous) {
        await api.delete(`/favorites/${id}`)
      } else {
        await api.post(`/favorites/${id}`)
      }
      return !previous
    } catch (e) {
      setFavorited(id, previous)
      console.error('Failed toggling favorite', e)
      return previous
    }
  }

  function remove(id: number) {
    setFavorited(id, false)
  }

  return {
    favoritedIds,
    loaded,
    hydrate,
    isFavorited,
    isFavoritedReactive,
    toggleFavorite,
    remove,
  }
}
