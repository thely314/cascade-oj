import { useToast } from 'vue-toastification';

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
