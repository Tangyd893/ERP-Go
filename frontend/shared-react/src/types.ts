export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface PageRequest {
  page: number;
  page_size: number;
}

export interface UserInfo {
  user_id: string;
  username: string;
  nickname: string;
  email: string;
  tenant_id: string;
  roles: string[];
  permissions: string[];
}
