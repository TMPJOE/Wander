<script setup lang="ts">
import { computed } from 'vue'
import { Calendar, Users, MapPin, CreditCard } from '@lucide/vue'

const props = defineProps<{
  booking: {
    id: number
    tour_title: string
    tour_location: string
    tour_image?: string
    schedule_start: string
    guest_count: number
    total_price: number
    status: string
  }
}>()

defineEmits<{
  cancel: [id: number]
}>()

const imageUrl = computed(() => {
  if (!props.booking.tour_image)
    return 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=400&h=300&fit=crop'
  try {
    const imgs = JSON.parse(props.booking.tour_image)
    return (
      imgs[0] || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=400&h=300&fit=crop'
    )
  } catch {
    return props.booking.tour_image
  }
})

const formattedDate = computed(() => {
  return new Date(props.booking.schedule_start).toLocaleDateString('es-MX', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
})

const formattedTime = computed(() => {
  return new Date(props.booking.schedule_start).toLocaleTimeString('es-MX', {
    hour: '2-digit',
    minute: '2-digit',
  })
})

const statusLabel: Record<string, string> = {
  pending: 'Pendiente',
  confirmed: 'Confirmada',
  cancelled: 'Cancelada',
  completed: 'Completada',
}

const statusClass: Record<string, string> = {
  pending: 'badge-warning',
  confirmed: 'badge-success',
  cancelled: 'badge-error',
  completed: 'badge-primary',
}
</script>

<template>
  <div class="booking-card card">
    <div class="booking-card__header">
      <span class="badge" :class="statusClass[booking.status] || 'badge-secondary'">
        {{ statusLabel[booking.status] || booking.status }}
      </span>
      <span class="booking-card__id">Reserva #{{ booking.id }}</span>
    </div>

    <div class="booking-card__body">
      <img :src="imageUrl" :alt="booking.tour_title" class="booking-card__image" />
      <div class="booking-card__info">
        <h3 class="booking-card__title">{{ booking.tour_title }}</h3>
        <p class="booking-card__location">
          <MapPin :size="14" />
          {{ booking.tour_location }}
        </p>

        <div class="booking-card__details">
          <div class="detail-item">
            <Calendar :size="14" class="detail-icon" />
            <span>{{ formattedDate }} • {{ formattedTime }}</span>
          </div>
          <div class="detail-item">
            <Users :size="14" class="detail-icon" />
            <span>{{ booking.guest_count }} persona{{ booking.guest_count > 1 ? 's' : '' }}</span>
          </div>
          <div class="detail-item">
            <CreditCard :size="14" class="detail-icon" />
            <span class="font-semibold">${{ booking.total_price.toLocaleString('es-MX') }}</span>
          </div>
        </div>
      </div>
    </div>

    <div
      class="booking-card__footer"
      v-if="booking.status === 'pending' || booking.status === 'confirmed'"
    >
      <button class="btn btn-outline btn-md" @click="$emit('cancel', booking.id)">
        Cancelar reserva
      </button>
    </div>
  </div>
</template>

<style scoped>
.booking-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.booking-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--color-border-light);
  padding: var(--spacing-3);
  margin-bottom: var(--spacing-1);
}

.booking-card__id {
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  font-weight: var(--font-weight-medium);
  margin-right: var(--spacing-2);
}

.booking-card__body {
  padding-top: var(--spacing-3);
  padding-bottom: var(--spacing-3);
  padding-left: var(--spacing-8);
  padding-right: var(--spacing-8);
  display: flex;
  gap: var(--spacing-4);
}

.booking-card__image {
  width: 90px;
  height: 90px;
  border-radius: var(--radius-md);
  object-fit: cover;
  flex-shrink: 0;
}

.booking-card__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.booking-card__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  line-height: var(--line-height-tight);
}

.booking-card__location {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  margin-bottom: var(--spacing-2);
}

.booking-card__details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.detail-icon {
  color: var(--color-text-light);
}

.font-semibold {
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.booking-card__footer {
  margin-top: var(--spacing-2);
  padding-right: var(--spacing-3);
  padding-bottom: var(--spacing-3);
  display: flex;
  justify-content: flex-end;
}
</style>
