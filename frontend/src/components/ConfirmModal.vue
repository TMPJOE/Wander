<script setup lang="ts">
import {
  isConfirmModalOpen,
  confirmModalTitle,
  confirmModalBody,
  confirmModalConfirmText,
  confirmModalCancelText,
  confirmModalVariant,
  resolveConfirmPromise,
} from '../composables/useConfirm'

const handleConfirm = () => {
  isConfirmModalOpen.value = false
  if (resolveConfirmPromise) resolveConfirmPromise(true)
}

const handleCancel = () => {
  isConfirmModalOpen.value = false
  if (resolveConfirmPromise) resolveConfirmPromise(false)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="isConfirmModalOpen" class="modal-overlay" @click.self="handleCancel">
        <div class="modal-card" role="dialog" aria-modal="true">
          <div class="modal-header">
            <h3 class="modal-title">{{ confirmModalTitle }}</h3>
          </div>
          <div class="modal-body">
            <p>{{ confirmModalBody }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="handleCancel">
              {{ confirmModalCancelText }}
            </button>
            <button
              class="btn"
              :class="confirmModalVariant === 'danger' ? 'btn-danger' : 'btn-primary'"
              @click="handleConfirm"
            >
              {{ confirmModalConfirmText }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9990;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.modal-card {
  background: var(--color-surface, #ffffff);
  border-radius: 16px;
  max-width: 440px;
  width: 100%;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.modal-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text, #0f172a);
}

.modal-body p {
  margin: 0;
  color: #475569;
  font-size: 0.95rem;
  line-height: 1.5;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.btn {
  padding: 0.625rem 1.25rem;
  border-radius: 10px;
  font-weight: 600;
  font-size: 0.875rem;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary {
  background: #f1f5f9;
  color: #475569;
}

.btn-secondary:hover {
  background: #e2e8f0;
}

.btn-primary {
  background: var(--color-primary, #6366f1);
  color: #ffffff;
}

.btn-primary:hover {
  filter: brightness(1.1);
}

.btn-danger {
  background: #ef4444;
  color: #ffffff;
}

.btn-danger:hover {
  background: #dc2626;
}

/* Modal Transition */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.25s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-active .modal-card {
  animation: modalPop 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalPop {
  0% {
    opacity: 0;
    transform: scale(0.92) translateY(10px);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
</style>
