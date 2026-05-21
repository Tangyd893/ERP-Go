import { defineStore } from "pinia";
import { ref } from "vue";
import apiClient from "@erp/shared";
import type { ApiResponse, PageData } from "@erp/shared";

interface Product {
  id: string;
  code: string;
  name: string;
  spu_name: string;
  barcode: string;
  weight: number;
  sale_price: number;
  currency: string;
  status: string;
}

interface SKUMappingRecord {
  id: string;
  sku_code: string;
  sku_name: string;
  platform_code: string;
  platform_sku: string;
  asin: string;
  fnsku: string;
  store: string;
}

export const useProductStore = defineStore("product", () => {
  const products = ref<Product[]>([]);
  const productTotal = ref(0);
  const skuMappings = ref<SKUMappingRecord[]>([]);
  const loading = ref(false);

  async function fetchProducts(page: number, pageSize: number) {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<Product>>>(
        "/products",
        { params: { page, page_size: pageSize } }
      );
      products.value = res.data.data.list;
      productTotal.value = res.data.data.total;
    } finally {
      loading.value = false;
    }
  }

  async function fetchSKUMappings() {
    loading.value = true;
    try {
      const res = await apiClient.get<ApiResponse<PageData<SKUMappingRecord>>>(
        "/products/sku-mappings"
      );
      skuMappings.value = res.data.data.list;
    } finally {
      loading.value = false;
    }
  }

  async function createSKUMapping(data: Partial<SKUMappingRecord>) {
    const res = await apiClient.post<ApiResponse<SKUMappingRecord>>(
      "/products/sku-mappings",
      data
    );
    return res.data.data;
  }

  return {
    products,
    productTotal,
    skuMappings,
    loading,
    fetchProducts,
    fetchSKUMappings,
    createSKUMapping,
  };
});
