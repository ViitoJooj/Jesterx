import { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useUser } from "./UserContext";

type AdminRouteProps = {
  children: ReactNode;
};

export function AdminRoute({ children }: AdminRouteProps) {
  const { user } = useUser();
  const location = useLocation();

  if (user === undefined) {
    return null;
  }

  if (user === null) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  if (user.role !== "platform_admin") {
    return <Navigate to="/404" replace />;
  }

  return <>{children}</>;
}
