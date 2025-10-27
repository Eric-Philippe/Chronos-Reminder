import { toast } from "sonner";
import type { ExternalToast } from "sonner";

interface UseToastOptions extends ExternalToast {
  description?: string;
}

export function useToast() {
  return {
    success: (message: string, options?: UseToastOptions) =>
      toast.success(message, options),
    error: (message: string, options?: UseToastOptions) =>
      toast.error(message, options),
    info: (message: string, options?: UseToastOptions) =>
      toast.info(message, options),
    warning: (message: string, options?: UseToastOptions) =>
      toast.warning(message, options),
    loading: (message: string, options?: UseToastOptions) =>
      toast.loading(message, options),
    promise: toast.promise,
  };
}
