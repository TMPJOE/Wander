<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useApi } from '../composables/useApi';
import { useToast } from '../composables/useToast';
import { ArrowLeft, Eye } from '@lucide/vue';
import TourForm from '../components/TourForm.vue';

const route = useRoute();
const router = useRouter();
const api = useApi();
const toast = useToast();

const tourId = computed(() => route.params.id as string);
const isEdit = computed(() => !!tourId.value);
const loading = ref(false);
const initialData = ref<any>(null);

// Holds the latest form state emitted by TourForm so we can preview it
const latestFormData = ref<any>(null);

onMounted(async () => {
  if (isEdit.value) {
    loading.value = true;
    try {
      const res = await api.get(`/tours/${tourId.value}`);
      initialData.value = res.data;
      latestFormData.value = res.data;
    } catch (e) {
      console.error(e);
      toast.error('Error cargando tour');
      router.back();
    } finally {
      loading.value = false;
    }
  }
});

function handleChange(data: any) {
  latestFormData.value = data;
}

async function handleSubmit(data: any) {
  loading.value = true;
  try {
    if (isEdit.value) {
      await api.put(`/tours/${tourId.value}`, data);
      toast.success('Tour actualizado exitosamente');
    } else {
      await api.post('/tours', data);
      toast.success('Tour creado exitosamente');
    }
    router.push('/guide/tours');
  } catch (e) {
    console.error(e);
    toast.error('Error al guardar el tour');
  } finally {
    loading.value = false;
  }
}

function handlePreview() {
  if (!latestFormData.value) {
    toast.info('Completa el formulario antes de previsualizar');
    return;
  }
  sessionStorage.setItem('wander_tour_preview', JSON.stringify(latestFormData.value));
  const returnPath = isEdit.value
    ? `/guide/tours/${tourId.value}/edit`
    : '/guide/tours/new';
  sessionStorage.setItem('wander_preview_return', returnPath);
  // Navigate to preview of existing tour or a synthetic one
  if (isEdit.value) {
    router.push(`/tours/${tourId.value}?preview=1`);
  } else {
    router.push(`/preview/tour?preview=1`);
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
      <button class="preview-btn" @click="handlePreview" title="Vista previa como viajero">
        <Eye :size="18" />
        <span>Previsualizar</span>
      </button>
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
        @change="handleChange"
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

.preview-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-1);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-full);
  background: var(--color-secondary-50);
  color: var(--color-secondary);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  border: 1px solid var(--color-secondary-100);
  transition: all var(--transition-fast);
  cursor: pointer;
}

.preview-btn:hover {
  background: var(--color-secondary-100);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
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
