<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useApi } from '../composables/useApi';
import { ArrowLeft } from '@lucide/vue';
import TourForm from '../components/TourForm.vue';

const route = useRoute();
const router = useRouter();
const api = useApi();

const tourId = computed(() => route.params.id as string);
const isEdit = computed(() => !!tourId.value);
const loading = ref(false);
const initialData = ref<any>(null);

onMounted(async () => {
  if (isEdit.value) {
    loading.value = true;
    try {
      const res = await api.get(`/tours/${tourId.value}`);
      initialData.value = res.data;
    } catch (e) {
      console.error(e);
      alert('Error cargando tour');
      router.back();
    } finally {
      loading.value = false;
    }
  }
});

async function handleSubmit(data: any) {
  loading.value = true;
  try {
    if (isEdit.value) {
      await api.put(`/tours/${tourId.value}`, data);
    } else {
      await api.post('/tours', data);
    }
    router.push('/guide/tours');
  } catch (e) {
    console.error(e);
    alert('Error al guardar el tour');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="page bg-surface min-h-screen">
    <header class="header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="20" />
      </button>
      <h1 class="title">{{ isEdit ? 'Editar Tour' : 'Crear Tour' }}</h1>
      <div style="width: 36px"></div>
    </header>

    <div class="container py-4">
      <div v-if="loading && isEdit && !initialData" class="flex justify-center p-8">
        Cargando...
      </div>
      <TourForm 
        v-else 
        :initial-data="initialData" 
        :loading="loading" 
        @submit="handleSubmit" 
        @cancel="router.back()" 
      />
    </div>
  </div>
</template>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  position: sticky;
  top: 0;
  z-index: 10;
}

.back-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  background: var(--color-background);
}

.title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
}

.bg-surface { background: var(--color-surface); }
.min-h-screen { min-height: 100vh; }
.py-4 { padding-top: var(--spacing-4); padding-bottom: var(--spacing-4); }
.flex { display: flex; }
.justify-center { justify-content: center; }
.p-8 { padding: var(--spacing-8); }
</style>
