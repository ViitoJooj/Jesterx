import "./styles/global.scss";
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import { AuthProvider } from "./hooks/AuthContext";
import { DEFAULT_WEBSITE_ID } from "./hooks/api";

import { Header } from "./components/header/Header";
import { Footer } from "./components/footer/Footer";
import { ProtectedRoute } from "./components/protectedRoute/ProtectedRoute";

const Home = React.lazy(() => import("./pages/home/Home").then((m) => ({ default: m.Home })));
const Register = React.lazy(() => import("./pages/register/Register").then((m) => ({ default: m.Register })));
const Login = React.lazy(() => import("./pages/login/Login").then((m) => ({ default: m.Login })));
const VerifyEmail = React.lazy(() => import("./pages/verify-email/VerifyEmail").then((m) => ({ default: m.VerifyEmail })));
const Plans = React.lazy(() => import("./pages/plans/Plans").then((m) => ({ default: m.Plans })));
const Pages = React.lazy(() => import("./pages/pages/Pages").then((m) => ({ default: m.Pages })));
const ElementorEditor = React.lazy(() => import("./pages/elementor-editor/ElementorEditor").then((m) => ({ default: m.ElementorEditor })));
const ReactEditor = React.lazy(() => import("./pages/react-editor/ReactEditor").then((m) => ({ default: m.ReactEditor })));
const SvelteEditor = React.lazy(() => import("./pages/svelte-editor/SvelteEditor").then((m) => ({ default: m.SvelteEditor })));
const PaymentSuccess = React.lazy(() => import("./pages/payment-success/PaymentSuccess").then((m) => ({ default: m.PaymentSuccess })));
const PaymentCancel = React.lazy(() => import("./pages/payment-cancel/PaymentCancel").then((m) => ({ default: m.PaymentCancel })));
const NotFound = React.lazy(() => import("./pages/not-found/NotFound").then((m) => ({ default: m.NotFound })));
const Docs = React.lazy(() => import("./pages/docs/Docs").then((m) => ({ default: m.Docs })));
const About = React.lazy(() => import("./pages/about/About").then((m) => ({ default: m.About })));
const Status = React.lazy(() => import("./pages/status/Status").then((m) => ({ default: m.Status })));
const Profile = React.lazy(() => import("./pages/profile/Profile").then((m) => ({ default: m.Profile })));
const AdminDashboard = React.lazy(() => import("./pages/admin/AdminDashboard").then((m) => ({ default: m.AdminDashboard })));
const StoreInfo = React.lazy(() => import("./pages/store-info/StoreInfo").then((m) => ({ default: m.StoreInfo })));
const Report = React.lazy(() => import("./pages/report/Report").then((m) => ({ default: m.Report })));

const websiteId = (import.meta.env.VITE_WEBSITE_ID as string | undefined) ?? DEFAULT_WEBSITE_ID;

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <AuthProvider websiteId={websiteId}>
        <Header />
        <React.Suspense fallback={<div style={{ padding: "1rem" }}>Carregando...</div>}>
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
        </React.Suspense>
        <Footer />
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>
);
