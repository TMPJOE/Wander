import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '../composables/useApi';

export const useFavoritesStore = defineStore('favorites', () => {
  const favorites = ref<any[]>([]);
  const loading = ref(false);
  const api = useApi();

  async function fetchFavorites() {
    loading.value = true;
    try {
      const response = await api.get('/favorites');
      favorites.value = response.data;
    } catch (e) {
      console.error('Failed to fetch favorites', e);
    } finally {
      loading.value = false;
    }
  }

  async function toggleFavorite(tourId: string, isFavorited: boolean) {
    try {
      if (isFavorited) {
        await api.delete(`/favorites/${tourId}`);
      } else {
        await api.post(`/favorites/${tourId}`);
      }
      await fetchFavorites();
    } catch (e) {
      console.error('Failed to toggle favorite', e);
    }
  }

  return { favorites, loading, fetchFavorites, toggleFavorite };
});
