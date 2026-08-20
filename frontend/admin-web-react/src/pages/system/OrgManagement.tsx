import { useQuery } from "@tanstack/react-query";
import {
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  ErrorState,
  Skeleton,
  apiClient,
} from "@erp/shared-react";
import { stubAction, unwrapList } from "@/lib/utils";

interface OrgNode {
  id: string;
  name: string;
  code: string;
  children?: OrgNode[];
}

function OrgTreeNode({ node, depth = 0 }: { node: OrgNode; depth?: number }) {
  return (
    <div>
      <div
        className="flex items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-muted/60"
        style={{ paddingLeft: `${depth * 12 + 8}px` }}
      >
        <span>{node.name}</span>
        <span className="text-xs text-muted-foreground">{node.code}</span>
      </div>
      {node.children?.map((child) => (
        <OrgTreeNode key={child.id} node={child} depth={depth + 1} />
      ))}
    </div>
  );
}

export default function OrgManagement() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["organizations"],
    queryFn: async () => {
      const res = await apiClient.get("/tenant/organizations");
      return unwrapList<OrgNode>(res.data?.data ?? res.data);
    },
  });

  if (isLoading) {
    return (
      <div className="space-y-3">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (error) {
    return <ErrorState title="加载组织架构失败" />;
  }

  if (!data?.length) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center gap-3 py-12 text-center">
          <p className="text-sm text-muted-foreground">暂无组织架构数据，请先创建组织</p>
          <Button disabled onClick={() => stubAction("新建组织")}>
            新建组织
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="grid gap-4 lg:grid-cols-[1fr_2fr]">
      <Card>
        <CardHeader>
          <CardTitle className="text-base">组织架构</CardTitle>
        </CardHeader>
        <CardContent className="space-y-1">
          {data.map((node) => (
            <OrgTreeNode key={node.id} node={node} />
          ))}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-base">部门/岗位</CardTitle>
          <Button size="sm" onClick={() => stubAction("新建部门")}>
            新建部门
          </Button>
        </CardHeader>
        <CardContent>
          <p className="py-8 text-center text-sm text-muted-foreground">请在左侧选择组织节点</p>
        </CardContent>
      </Card>
    </div>
  );
}
