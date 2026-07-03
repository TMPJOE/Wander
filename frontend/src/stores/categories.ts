import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '../composables/useApi';

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<any[]>([]);
  const loading = ref(false);
  const api = useApi();

  async function fetchCategories() {
    loading.value = true;
    try {
      const response = await api.get('/categories');
      categories.value = response.data;
    } catch (e) {
      console.error('Failed to fetch categories', e);
    } finally {
      loading.value = false;
    }
  }

  return { categories, loading, fetchCategories };
});
