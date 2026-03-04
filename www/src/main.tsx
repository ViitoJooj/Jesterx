import "./styles/global.scss";
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import { AuthProvider } from "./hooks/AuthContext";

import { Header } from "./components/header/Header";
import { Home } from "./pages/home/Home";
import { Register } from "./pages/register/Register";
import { Login } from "./pages/login/Login";
import { Plans } from "./pages/plans/Plans";
import { PaymentSuccess } from "./pages/payment-success/PaymentSuccess";
import { PaymentCancel } from "./pages/payment-cancel/PaymentCancel";
import { NotFound } from "./pages/not-found/NotFound";

const websiteId = import.meta.env.VITE_WEBSITE_ID as string;

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <AuthProvider websiteId={websiteId}>
        <Header />

        <Routes>
          <Route path="*" element={<NotFound />} />
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />

          <Route path="/plans" element={<Plans />} />
          <Route path="/payment-cancel" element={<PaymentCancel />} />
          <Route path="/payment-success" element={<PaymentSuccess />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>
);