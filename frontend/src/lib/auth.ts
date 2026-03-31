export interface AuthUser {
  id: string;
  email: string;
  role: string;
}

export function getAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("auth_token");
}

export function getAuthUser(): AuthUser | null {
  if (typeof window === "undefined") return null;
  const userStr = localStorage.getItem("auth_user");
  if (!userStr) return null;
  try {
    return JSON.parse(userStr) as AuthUser;
  } catch {
    return null;
  }
}

export function isAuthenticated(): boolean {
  return !!getAuthToken();
}

export function isAdmin(): boolean {
  const user = getAuthUser();
  return user?.role === "admin" || user?.role === "manager";
}
