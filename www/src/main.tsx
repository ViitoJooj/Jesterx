import "./styles/global.scss";
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import { AuthProvider } from "./hooks/AuthContext";
import { DEFAULT_WEBSITE_ID } from "./hooks/api";

import { Header } from "./components/header/Header";
import { Footer } from "./components/footer/Footer";
import { ProtectedRoute } from "./components/protectedRoute/ProtectedRoute";
import { Home } from "./pages/home/Home";
import { Register } from "./pages/register/Register";
import { Login } from "./pages/login/Login";
import { VerifyEmail } from "./pages/verify-email/VerifyEmail";
import { Plans } from "./pages/plans/Plans";
import { Pages } from "./pages/pages/Pages";
import { ElementorEditor } from "./pages/elementor-editor/ElementorEditor";
import { ReactEditor } from "./pages/react-editor/ReactEditor";
import { SvelteEditor } from "./pages/svelte-editor/SvelteEditor";
import { PaymentSuccess } from "./pages/payment-success/PaymentSuccess";
import { PaymentCancel } from "./pages/payment-cancel/PaymentCancel";
import { NotFound } from "./pages/not-found/NotFound";
import { Docs } from "./pages/docs/Docs";
import { About } from "./pages/about/About";
import { Status } from "./pages/status/Status";

import { Profile } from "./pages/profile/Profile";
import { AdminDashboard } from "./pages/admin/AdminDashboard";
import { StoreInfo } from "./pages/store-info/StoreInfo";
import { Report } from "./pages/report/Report";

const websiteId = (import.meta.env.VITE_WEBSITE_ID as string | undefined) ?? DEFAULT_WEBSITE_ID;

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <AuthProvider websiteId={websiteId}>
        <Header />
        <Routes>
          <Route path="*" element={<NotFound />} />
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<Register />} />
          <Route path="/verify-email" element={<VerifyEmail />} />
          <Route path="/login" element={<Login />} />
          <Route path="/plans" element={<Plans />} />
          <Route path="/pages" element={<ProtectedRoute><Pages /></ProtectedRoute>} />
          <Route path="/pages/:siteId/elementor" element={<ProtectedRoute><ElementorEditor /></ProtectedRoute>} />
          <Route path="/pages/:siteId/react" element={<ProtectedRoute><ReactEditor /></ProtectedRoute>} />
          <Route path="/pages/:siteId/svelte" element={<ProtectedRoute><SvelteEditor /></ProtectedRoute>} />
          <Route path="/payment-cancel" element={<ProtectedRoute><PaymentCancel /></ProtectedRoute>} />
          <Route path="/payment-success" element={<ProtectedRoute><PaymentSuccess /></ProtectedRoute>} />
          <Route path="/profile" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
          <Route path="/docs" element={<Docs />} />
          <Route path="/about" element={<About />} />
          <Route path="/status" element={<Status />} />
          <Route path="/admin" element={<ProtectedRoute><AdminDashboard /></ProtectedRoute>} />
          <Route path="/store/:siteId" element={<StoreInfo />} />
          <Route path="/report" element={<Report />} />
        </Routes>
        <Footer />
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>
);
