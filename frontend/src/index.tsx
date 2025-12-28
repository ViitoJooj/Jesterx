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
import { ThemeStore } from "./pages/ThemeStore";
import { ThemePreviewPage } from "./pages/ThemePreview";

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
        <Route path="/themes" element={<ThemeStore />} />
        <Route path="/themes/:themeId" element={<ThemePreviewPage />} />
        <Route path="/payment-success" element={<ProtectedRoute> <PaymentSuccess /> </ProtectedRoute>} />
        <Route path="/payment-cancel" element={<ProtectedRoute> <PaymentCancel /> </ProtectedRoute>} />
        <Route path="/pages/:pageId" element={<ProtectedRoute> <PagePreview /> </ProtectedRoute>} />
        <Route path="/pages/:pageId/edit" element={<ProtectedRoute> <PageEdit /> </ProtectedRoute>} />
      </Routes>
      <Footer />
      <SpeedInsights />
    </UserProvider>
  </BrowserRouter>
);
