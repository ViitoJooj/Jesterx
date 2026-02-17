import { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useUser } from "./UserContext";

type ProtectedRouteProps = {
  children: ReactNode;
};

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { user } = useUser();
  const location = useLocation();

  if (user === undefined) {
    return null;
  }

  if (user === null) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  return <>{children}</>;
}