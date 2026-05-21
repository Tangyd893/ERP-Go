// API 通用响应类型
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

// 分页数据
export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// 分页请求参数
export interface PageRequest {
  page: number;
  page_size: number;
}

// 用户信息
export interface UserInfo {
  user_id: string;
  username: string;
  nickname: string;
  email: string;
  tenant_id: string;
  roles: string[];
  permissions: string[];
}
