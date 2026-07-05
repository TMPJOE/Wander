import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'explore',
      component: () => import('../views/ExploreView.vue'),
    },
    {
      path: '/tours/:id',
      name: 'tour-detail',
      component: () => import('../views/TourDetailView.vue'),
    },
    {
      path: '/tours/:id/book',
      name: 'book-tour',
      component: () => import('../views/BookingView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/checkout/:bookingId',
      name: 'checkout',
      component: () => import('../views/CheckoutView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/booking-success/:id',
      name: 'booking-success',
      component: () => import('../views/BookingSuccessView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/bookings',
      name: 'my-bookings',
      component: () => import('../views/MyBookingsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/messages',
      name: 'messages',
      component: () => import('../views/MessagesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/messages/:userId',
      name: 'chat',
      component: () => import('../views/ChatView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
    },
    {
      path: '/guide/dashboard',
      name: 'guide-dashboard',
      component: () => import('../views/GuideDashboardView.vue'),
      meta: { requiresAuth: true, requiresGuide: true },
    },
    {
      path: '/guide/tours',
      name: 'guide-tours',
      component: () => import('../views/GuideToursView.vue'),
      meta: { requiresAuth: true, requiresGuide: true },
    },
    {
      path: '/guide/tours/new',
      name: 'tour-form-new',
      component: () => import('../views/TourFormView.vue'),
      meta: { requiresAuth: true, requiresGuide: true },
    },
    {
      path: '/guide/tours/:id/edit',
      name: 'tour-form-edit',
      component: () => import('../views/TourFormView.vue'),
      meta: { requiresAuth: true, requiresGuide: true },
    },
    {
      path: '/guide/tours/:id/schedules',
      name: 'schedule-manager',
      component: () => import('../views/ScheduleManagerView.vue'),
      meta: { requiresAuth: true, requiresGuide: true },
    },
    {
      path: '/guide/bookings',
      name: 'guide-bookings',
      component: () => import('../views/GuideBookingsView.vue'),
      meta: { requiresAuth: true, requiresGuide: true },
    },
    {
      path: '/category/:slug',
      name: 'category',
      component: () => import('../views/CategoryView.vue'),
    },
  ],
})

let isInitialLoad = true

router.beforeEach(async (to, _from) => {
  const authStore = useAuthStore()

  // On first load, if we have a token but no user object, fetch the user
  if (isInitialLoad && authStore.token && !authStore.user) {
    isInitialLoad = false
    await authStore.fetchMe()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresGuide && !authStore.isGuide) {
    return { name: 'explore' }
  }

  // allow navigation
  return true
})

export default router
