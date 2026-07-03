<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useApi } from '../composables/useApi';
import { ArrowLeft, Plus, Edit2, Calendar } from '@lucide/vue';

const router = useRouter();
const api = useApi();
const tours = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const res = await api.get('/guide/tours');
    tours.value = res.data || [];
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
});

function getImageUrl(images: string) {
  try {
    const parsed = JSON.parse(images);
    return parsed[0] || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=400&h=300&fit=crop';
  } catch {
    return images || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=400&h=300&fit=crop';
  }
}
</script>

<template>
  <div class="page">
    <header class="header">
      <div class="flex items-center gap-3">
        <button class="back-btn" @click="router.push('/guide/dashboard')">
          <ArrowLeft :size="20" />
        </button>
        <h1 class="title">Mis Tours</h1>
      </div>
      <button class="btn btn-primary btn-sm" @click="router.push('/guide/tours/new')">
        <Plus :size="16" /> Nuevo
      </button>
    </header>

    <div class="container py-4">
      <div v-if="loading" class="flex flex-col gap-4">
        <div v-for="i in 3" :key="i" class="skeleton h-24 rounded-lg"></div>
      </div>
      
      <div v-else-if="tours.length" class="flex flex-col gap-4">
        <div v-for="tour in tours" :key="tour.id" class="tour-item card p-3 flex gap-3">
          <img :src="getImageUrl(tour.images)" class="tour-img" />
          <div class="tour-info">
            <h3 class="font-semibold text-sm line-clamp-1">{{ tour.title }}</h3>
            <p class="text-xs text-secondary mt-1">${{ tour.price_per_person }} MXN • {{ tour.duration_minutes }} min</p>
            
            <div class="flex gap-2 mt-auto pt-2">
              <button class="btn btn-outline btn-xs flex-1" @click="router.push(`/guide/tours/${tour.id}/edit`)">
                <Edit2 :size="14" /> Editar
              </button>
              <button class="btn btn-primary btn-xs flex-1" @click="router.push(`/guide/tours/${tour.id}/schedules`)">
                <Calendar :size="14" /> Horarios
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="text-center py-12">
        <p class="text-secondary mb-4">No has creado ningún tour todavía.</p>
        <button class="btn btn-primary" @click="router.push('/guide/tours/new')">
          Crear mi primer tour
        </button>
      </div>
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

.tour-item {
  flex-direction: row;
}

.tour-img {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-md);
  object-fit: cover;
  flex-shrink: 0;
}

.tour-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.flex { display: flex; }
.flex-col { flex-direction: column; }
.gap-3 { gap: var(--spacing-3); }
.gap-4 { gap: var(--spacing-4); }
.gap-2 { gap: var(--spacing-2); }
.items-center { align-items: center; }
.py-4 { padding-top: var(--spacing-4); padding-bottom: var(--spacing-4); }
.p-3 { padding: var(--spacing-3); }
.h-24 { height: 6rem; }
.rounded-lg { border-radius: var(--radius-lg); }
.font-semibold { font-weight: var(--font-weight-semibold); }
.text-sm { font-size: var(--font-size-sm); }
.text-xs { font-size: var(--font-size-xs); }
.text-secondary { color: var(--color-text-secondary); }
.mt-1 { margin-top: 2px; }
.mt-auto { margin-top: auto; }
.pt-2 { padding-top: var(--spacing-2); }
.flex-1 { flex: 1; }
.text-center { text-align: center; }
.py-12 { padding-top: 3rem; padding-bottom: 3rem; }
.mb-4 { margin-bottom: var(--spacing-4); }
</style>
