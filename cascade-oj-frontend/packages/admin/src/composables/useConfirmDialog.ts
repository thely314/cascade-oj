import { useToast } from 'vue-toastification';

/**
 * Returns a confirm dialog function that uses toast instead of native confirm().
 * Usage: const confirm = useConfirmDialog();
 *        const ok = await confirm('确定要删除吗？');
 */
export function useConfirmDialog() {
  const toast = useToast();

  return (message: string): Promise<boolean> => {
    return new Promise((resolve) => {
      let resolved = false;
      const id = toast.info(message, {
        timeout: 8000,
        closeOnClick: true,
        onClick: () => {
          if (!resolved) {
            resolved = true;
            toast.dismiss(id);
            resolve(true);
          }
        },
        onClose: () => {
          if (!resolved) {
            resolved = true;
            resolve(false);
          }
        },
      });
    });
  };
}
