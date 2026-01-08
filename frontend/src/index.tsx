import "./styles/global.scss";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Header } from "./components/Header";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { UserProvider } from "./config/UserContext";
import { ProtectedRoute } from "./config/ProtectedRoute";
import { Footer } from "./components/Footer";
import { NotFound } from "./pages/Not-found";
import { HomePage } from "./pages/HomePage";
import { SpeedInsights } from "@vercel/speed-insights/react";
import { Pages } from "./pages/Pages";
import { Products } from "./pages/Products";
import { Pricing } from "./pages/Pricing";
import { PaymentSuccess } from "./pages/PaymentSuccess";
import { PaymentCancel } from "./pages/PaymentCancel";
import { PagePreview } from "./pages/PagePreview";
import { PageEdit } from "./pages/PageEdit";
import { AdminRoute } from "./config/AdminRoute";
import { AdminDashboard } from "./pages/AdminDashboard";
import { ThemeStore } from "./pages/ThemeStore";
import { PublicPage } from "./pages/PublicPage";
import { ThemeProduct } from "./pages/ThemeProduct";

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);

root.render(
  <BrowserRouter>
    <UserProvider>
      <Header />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/404" element={<NotFound />} />
        <Route path="*" element={<NotFound />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/pages" element={<ProtectedRoute> <Pages /> </ProtectedRoute>} />
        <Route path="/products" element={<ProtectedRoute> <Products /> </ProtectedRoute>} />
        <Route path="/pricing" element={<Pricing />} />
        <Route path="/temas" element={<ThemeStore />} />
        <Route path="/payment-success" element={<ProtectedRoute> <PaymentSuccess /> </ProtectedRoute>} />
        <Route path="/payment-cancel" element={<ProtectedRoute> <PaymentCancel /> </ProtectedRoute>} />
        <Route path="/pages/:pageId" element={<ProtectedRoute> <PagePreview /> </ProtectedRoute>} />
        <Route path="/pages/:pageId/edit" element={<ProtectedRoute> <PageEdit /> </ProtectedRoute>} />
        <Route path="/public/:tenant/:pageId" element={<PublicPage />} />
        <Route path="/admin" element={<AdminRoute> <AdminDashboard /> </AdminRoute>} />
        <Route path="/themes/:slug" element={<ThemeProduct />} />
      </Routes>
      <Footer />
      <SpeedInsights />
    </UserProvider>
  </BrowserRouter>
);
