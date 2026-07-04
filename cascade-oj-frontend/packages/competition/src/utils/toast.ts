import type { ToastInterface } from 'vue-toastification';

let _toast: ToastInterface | null = null;

export function initToast(toast: ToastInterface): void {
  _toast = toast;
}

export function getToast(): ToastInterface {
  if (!_toast) {
    throw new Error(
      'Toast not initialized. Call initToast() after app.use(Toast).mount().',
    );
  }
  return _toast;
}
