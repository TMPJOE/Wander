import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '../composables/useApi';

export const useToursStore = defineStore('tours', () => {
  const tours = ref<any[]>([]);
  const currentTour = ref<any>(null);
  const loading = ref(false);
  const api = useApi();

  async function fetchTours(params?: any) {
    loading.value = true;
    try {
      const response = await api.get('/tours', { params });
      tours.value = response.data;
    } catch (e) {
      console.error('Failed to fetch tours', e);
    } finally {
      loading.value = false;
    }
  }

  async function fetchTourById(id: string) {
    loading.value = true;
    try {
      const response = await api.get(`/tours/${id}`);
      currentTour.value = response.data;
    } catch (e) {
      console.error('Failed to fetch tour', e);
    } finally {
      loading.value = false;
    }
  }

  return { tours, currentTour, loading, fetchTours, fetchTourById };
});
