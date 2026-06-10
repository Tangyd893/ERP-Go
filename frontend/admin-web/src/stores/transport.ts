import { defineStore } from "pinia";
import { ref } from "vue";
import { apiClient } from "@erp/shared";

export interface Carrier {
  id: string;
  name: string;
  code: string;
  contactPerson: string;
  contactPhone: string;
  status: string;
}

export interface Shipment {
  id: string;
  orderId: string;
  orderNo: string;
  outboundId: string;
  trackingNo: string;
  carrierId: string;
  carrierName: string;
  status: string;
  labelUrl: string;
  shippingCost: number;
  estimatedDelivery: string;
  actualDelivery: string;
  createdAt: string;
}

export const useTransportStore = defineStore("transport", () => {
  const carriers = ref<Carrier[]>([]);
  const shipments = ref<Shipment[]>([]);
  const loading = ref(false);
  const total = ref(0);

  async function fetchCarriers() {
    loading.value = true;
    try {
      const res = await apiClient.get("/transport/carriers");
      carriers.value = res.data?.data ?? res.data ?? [];
    } finally {
      loading.value = false;
    }
  }

  async function fetchShipments(params?: Record<string, any>) {
    loading.value = true;
    try {
      const res = await apiClient.get("/transport/shipments", { params });
      const data = res.data?.data ?? res.data;
      shipments.value = data.list ?? data.items ?? [];
      total.value = data.total ?? 0;
    } finally {
      loading.value = false;
    }
  }

  return {
    carriers,
    shipments,
    loading,
    total,
    fetchCarriers,
    fetchShipments,
    generateLabel,
  };
});

async function generateLabel(shipmentId: string) {
  const res = await apiClient.post(`/transport/shipments/${shipmentId}/label`);
  return res.data?.data ?? res.data;
}
