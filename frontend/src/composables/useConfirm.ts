import { ref } from 'vue'

export interface ConfirmOptions {
  title: string
  body: string
  confirmLabel?: string
  cancelLabel?: string
  confirmVariant?: 'danger' | 'primary'
}

export const isConfirmModalOpen = ref(false)
export const confirmModalTitle = ref('')
export const confirmModalBody = ref('')
export const confirmModalConfirmText = ref('Confirm')
export const confirmModalCancelText = ref('Cancel')
export const confirmModalVariant = ref<'danger' | 'primary'>('primary')
export let resolveConfirmPromise: ((value: boolean) => void) | null = null

export function useConfirm() {
  const confirm = (options: ConfirmOptions): Promise<boolean> => {
    confirmModalTitle.value = options.title
    confirmModalBody.value = options.body
    confirmModalConfirmText.value = options.confirmLabel || 'Confirm'
    confirmModalCancelText.value = options.cancelLabel || 'Cancel'
    confirmModalVariant.value = options.confirmVariant || 'primary'
    isConfirmModalOpen.value = true

    return new Promise<boolean>((resolve) => {
      resolveConfirmPromise = resolve
    })
  }

  return { confirm }
}
