import { ref, type Ref } from "vue";

/**
 * 统一的 API 调用状态管理 composable。
 * 替代各 Store 中手动重复的 loading/error 模式。
 *
 * @example
 * const { data, loading, error, execute } = useApiState<Order[]>();
 * await execute(() => apiClient.get("/orders"));
 */
export function useApiState<T>() {
  const data: Ref<T | null> = ref(null);
  const loading = ref(false);
  const error = ref("");

  async function execute(fn: () => Promise<T>): Promise<T | null> {
    loading.value = true;
    error.value = "";
    try {
      const result = await fn();
      data.value = result;
      return result;
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } })?.response?.data?.message ||
        (e as { message?: string })?.message ||
        "请求失败，请重试";
      error.value = msg;
      return null;
    } finally {
      loading.value = false;
    }
  }

  /** 重置状态 */
  function reset() {
    data.value = null;
    loading.value = false;
    error.value = "";
  }

  return { data, loading, error, execute, reset };
}

/**
 * 分页 API 调用状态管理。
 */
export function useApiPage<T>() {
  const list: Ref<T[]> = ref([]);
  const total = ref(0);
  const loading = ref(false);
  const error = ref("");

  async function execute(
    fn: () => Promise<{ data?: { list?: T[]; total?: number } }>
  ): Promise<void> {
    loading.value = true;
    error.value = "";
    try {
      const res = await fn();
      const pageData = (res as { data?: { list?: T[]; total?: number } }).data;
      list.value = pageData?.list ?? [];
      total.value = pageData?.total ?? 0;
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } })?.response?.data?.message ||
        (e as { message?: string })?.message ||
        "加载失败，请重试";
      error.value = msg;
    } finally {
      loading.value = false;
    }
  }

  function reset() {
    list.value = [];
    total.value = 0;
    loading.value = false;
    error.value = "";
  }

  return { list, total, loading, error, execute, reset };
}
