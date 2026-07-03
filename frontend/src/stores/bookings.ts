import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '../composables/useApi';

export const useBookingsStore = defineStore('bookings', () => {
  const bookings = ref<any[]>([]);
  const loading = ref(false);
  const api = useApi();

  async function fetchMyBookings() {
    loading.value = true;
    try {
      const response = await api.get('/bookings');
      bookings.value = response.data;
    } catch (e) {
      console.error('Failed to fetch bookings', e);
    } finally {
      loading.value = false;
    }
  }

  async function createBooking(data: any) {
    const response = await api.post('/bookings', data);
    return response.data;
  }

  return { bookings, loading, fetchMyBookings, createBooking };
});
