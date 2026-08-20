import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Badge,
  Button,
  Input,
  Label,
  ProTable,
  type ProColumn,
  apiClient,
  type ApiResponse,
  type PageData,
} from "@erp/shared-react";
import { PageCard } from "@/components/PageCard";
import { stubAction } from "@/lib/utils";

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

export default function ProductList() {
  const [page, setPage] = useState(1);
  const pageSize = 20;

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["products", page],
    queryFn: async () => {
      const res = await apiClient.get<ApiResponse<PageData<Product>>>("/products", {
        params: { page, page_size: pageSize },
      });
      return res.data.data;
    },
  });

  const columns: ProColumn<Product>[] = [
    { key: "code", title: "SKU编码", width: 150 },
    { key: "name", title: "SKU名称" },
    { key: "spu_name", title: "所属SPU", width: 120 },
    { key: "barcode", title: "条码", width: 150 },
    {
      key: "weight",
      title: "重量(kg)",
      width: 100,
      align: "right",
      render: (row) => row.weight,
    },
    {
      key: "sale_price",
      title: "售价",
      width: 120,
      align: "right",
      render: (row) => `${row.currency} ${row.sale_price.toFixed(2)}`,
    },
    {
      key: "status",
      title: "状态",
      width: 100,
      render: (row) => (
        <Badge variant={row.status === "active" ? "success" : "destructive"}>
          {row.status === "active" ? "启售" : "停售"}
        </Badge>
      ),
    },
  ];

  return (
    <PageCard
      title="商品列表 (SKU)"
      actions={
        <>
          <Button size="sm" variant="outline" disabled onClick={() => stubAction("新建 SPU")}>
            新建 SPU
          </Button>
          <Button size="sm" disabled onClick={() => stubAction("新建 SKU")}>
            新建 SKU
          </Button>
        </>
      }
    >
      <div className="mb-4 grid gap-3 md:grid-cols-3">
        <div className="space-y-1">
          <Label>SKU编码</Label>
          <Input placeholder="搜索编码" disabled />
        </div>
        <div className="space-y-1">
          <Label>商品名称</Label>
          <Input placeholder="搜索名称" disabled />
        </div>
        <div className="space-y-1">
          <Label>条码</Label>
          <Input placeholder="条码搜索" disabled />
        </div>
      </div>

      <ProTable
        columns={columns}
        data={data?.list ?? []}
        loading={isLoading}
        error={error ? "加载商品失败" : undefined}
        total={data?.total ?? 0}
        page={page}
        pageSize={pageSize}
        rowKey={(row) => row.id}
        onPageChange={setPage}
        onRefresh={() => refetch()}
        actions={() => (
          <>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("详情")}>
              详情
            </Button>
            <Button size="sm" variant="outline" disabled onClick={() => stubAction("编辑")}>
              编辑
            </Button>
          </>
        )}
      />
    </PageCard>
  );
}
