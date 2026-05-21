import { Navigate, useLocation } from "react-router-dom";
import { useAuthContext } from "../../hooks/AuthContext";

export function ProtectedRoute({ children }) {
  const { loading, isAuthenticated } = useAuthContext();
  const location = useLocation();

  if (loading) {
    return <main style={{ minHeight: "70vh", display: "grid", placeItems: "center" }}>Carregando...</main>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname + location.search }} />;
  }

  return <>{children}</>;
}
