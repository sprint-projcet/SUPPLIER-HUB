const SupplierDashboard = (() => {
  const NAV_BASE_CLASS =
    "w-full flex items-center gap-3 px-4 py-3.5 rounded-xl transition-all duration-200 group";
  const NAV_ACTIVE_CLASS = "bg-emerald-600 text-white shadow-xl shadow-emerald-900/40";
  const NAV_IDLE_CLASS = "text-slate-400 hover:bg-slate-800 hover:text-white";

  function getSession() {
    return typeof getStoredUserSession === "function"
      ? getStoredUserSession()
      : null;
  }

  function requireSupplier(redirectUrl = "../Login/login.html") {
    const user = checkAuth(redirectUrl);
    if (!user) return null;
    if (user.role !== "supplier") {
      window.location.href = redirectUrl;
      return null;
    }
    return user;
  }

  function getDisplayName(user) {
    return (
      user.business_name ||
      user.name ||
      user.email ||
      "Supplier"
    );
  }

  function setText(selector, value) {
    document.querySelectorAll(selector).forEach((element) => {
      element.textContent = value;
    });
  }

  function setUserDisplay(user) {
    const displayName = getDisplayName(user);
    setText("#user-display-name", displayName);
    setText("[data-supplier-name]", displayName);
    setText("[data-supplier-email]", user.email || "-");
    document.querySelectorAll("[data-supplier-name-display]").forEach((element) => {
      element.textContent = displayName;
    });
  }

  function setCurrentDate() {
    const dateEl = document.getElementById("current-date");
    if (!dateEl) return;

    dateEl.textContent = new Date().toLocaleDateString("id-ID", {
      day: "numeric",
      month: "long",
      year: "numeric",
    });
  }

  function setActiveNavigation() {
    const currentPage = window.location.pathname.split("/").pop() || "supplier.html";

    document.querySelectorAll("#sidebar nav a[href]").forEach((link) => {
      const isActive = link.getAttribute("href") === currentPage;
      link.className = `${NAV_BASE_CLASS} ${isActive ? NAV_ACTIVE_CLASS : NAV_IDLE_CLASS}`;

      link.querySelectorAll("i").forEach((icon) => {
        icon.className = isActive ? "text-white" : "group-hover:text-emerald-400";
      });
    });
  }

  function toggleSidebar() {
    const sidebar = document.getElementById("sidebar");
    if (sidebar) sidebar.classList.toggle("-translate-x-full");
  }

  async function apiFetch(path, options = {}) {
    const session = getSession();
    const headers = {
      ...(options.headers || {}),
    };

    if (!(options.body instanceof FormData)) {
      headers["Content-Type"] = headers["Content-Type"] || "application/json";
    }

    if (session && session.token) {
      headers.Authorization = `Bearer ${session.token}`;
    }

    const response = await fetch(buildApiUrl(path), {
      ...options,
      headers,
    });
    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      const error = new Error(data.error || "Request gagal diproses");
      error.status = response.status;
      error.data = data;
      throw error;
    }

    return data;
  }

  function getProfile() {
    return apiFetch("/api/supplier/profile").then((result) => result.data || result);
  }

  function updateProfile(payload) {
    return apiFetch("/api/supplier/profile", {
      method: "PUT",
      body: JSON.stringify(payload),
    }).then((result) => result.data || result);
  }

  function getStats() {
    return apiFetch("/api/supplier/stats");
  }

  function getProducts() {
    return apiFetch("/api/supplier/products").then((result) =>
      Array.isArray(result) ? result : result.data || [],
    );
  }

  function getOrders(status = "all") {
    const query = status && status !== "all" ? `?status=${encodeURIComponent(status)}` : "";
    return apiFetch(`/api/supplier/orders${query}`).then((result) => result.data || []);
  }

  function updateOrderStatus(orderID, status) {
    const normalized = String(status || "").toLowerCase();
    if (normalized === "confirm" || normalized === "reject") {
      return apiFetch(`/supplierhub/konfirmasi_stok/${encodeURIComponent(orderID)}`, {
        method: "PUT",
        body: JSON.stringify({ action: normalized }),
      });
    }

    return apiFetch(`/api/supplier/orders/${encodeURIComponent(orderID)}`, {
      method: "PUT",
      body: JSON.stringify({ status }),
    });
  }

  function getNotifications(options = {}) {
    const params = new URLSearchParams();
    if (options.unreadOnly) params.set("unread_only", "true");
    if (options.readStatus) params.set("read_status", options.readStatus);
    if (options.limit) params.set("limit", options.limit);
    const query = params.toString() ? `?${params.toString()}` : "";
    return apiFetch(`/api/supplier/notifications${query}`).then(
      (result) => result.data || [],
    );
  }

  function getUnreadNotifications() {
    return getNotifications({ unreadOnly: true });
  }

  function markNotificationRead(notificationID) {
    return apiFetch(`/api/supplier/notifications/${encodeURIComponent(notificationID)}/read`, {
      method: "PUT",
      body: JSON.stringify({}),
    });
  }

  function formatRupiah(value) {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(Number(value) || 0);
  }

  function formatNumber(value) {
    return new Intl.NumberFormat("id-ID").format(Number(value) || 0);
  }

  function escapeHTML(value) {
    return String(value ?? "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  }

  function shortID(value) {
    const id = String(value || "");
    if (!id) return "-";
    return `ORD-${id.slice(0, 8).toUpperCase()}`;
  }

  function statusMeta(status) {
    const normalized = String(status || "").toLowerCase();
    const map = {
      pending: {
        label: "PERLU DIPROSES",
        className: "bg-orange-50 text-orange-600",
        nextStatus: "confirm",
        actionLabel: "Konfirmasi",
      },
      pending_supplier_confirmation: {
        label: "MENUNGGU KONFIRMASI",
        className: "bg-orange-50 text-orange-600",
        nextStatus: "confirm",
        actionLabel: "Konfirmasi",
      },
      supplier_confirmed: {
        label: "STOK DIKONFIRMASI",
        className: "bg-teal-50 text-teal-600",
      },
      payment_pending: {
        label: "MENUNGGU BAYAR",
        className: "bg-amber-50 text-amber-600",
      },
      payment_request_failed: {
        label: "PAYMENT GAGAL DIBUAT",
        className: "bg-red-50 text-red-600",
      },
      payment_failed: {
        label: "PEMBAYARAN GAGAL",
        className: "bg-red-50 text-red-600",
      },
      rejected_by_supplier: {
        label: "DITOLAK SUPPLIER",
        className: "bg-red-50 text-red-600",
      },
      stock_unavailable: {
        label: "STOK TIDAK CUKUP",
        className: "bg-red-50 text-red-600",
      },
      paid: {
        label: "SIAP DIPROSES",
        className: "bg-emerald-50 text-emerald-600",
        nextStatus: "processing",
        actionLabel: "Proses",
      },
      processing: {
        label: "DIPROSES",
        className: "bg-blue-50 text-blue-600",
        nextStatus: "shipped",
        actionLabel: "Kirim",
      },
      shipped: {
        label: "DIKIRIM",
        className: "bg-indigo-50 text-indigo-600",
        nextStatus: "completed",
        actionLabel: "Selesaikan",
      },
      shipment_created: {
        label: "PENGIRIMAN DIBUAT",
        className: "bg-indigo-50 text-indigo-600",
        nextStatus: "completed",
        actionLabel: "Selesaikan",
      },
      completed: {
        label: "SELESAI",
        className: "bg-emerald-50 text-emerald-600",
      },
      cancelled: {
        label: "BATAL",
        className: "bg-red-50 text-red-600",
      },
    };

    return map[normalized] || {
      label: normalized ? normalized.toUpperCase() : "-",
      className: "bg-slate-100 text-slate-600",
    };
  }

  function notify(type, message) {
    if (typeof window.showGlobalToast === "function") {
      window.showGlobalToast(type, message);
      return;
    }

    if (typeof window.showToast === "function") {
      window.showToast(type, message);
    }
  }

  function updateStoredProfile(profile) {
    const current = getSession();
    if (!current || !profile) return;

    saveUserSession({
      ...current,
      name: profile.business_name || current.name,
      business_name: profile.business_name || current.business_name || "",
      email: profile.email || current.email || "",
      address: profile.address || "",
      category: profile.category || "",
      region: profile.region || "",
      status: profile.status || "",
    });
  }

  function renderNotificationPanel(notifications) {
    document.getElementById("supplier-notification-panel")?.remove();
    if (!Array.isArray(notifications) || notifications.length === 0) return;

    const panel = document.createElement("div");
    panel.id = "supplier-notification-panel";
    panel.className =
      "fixed right-5 top-24 z-[120] w-[calc(100%-2.5rem)] max-w-md space-y-3";

    notifications.slice(0, 5).forEach((notification) => {
      const card = document.createElement("div");
      card.className =
        "rounded-2xl border border-orange-100 bg-white p-4 shadow-2xl shadow-slate-900/10";

      const title = escapeHTML(notification.title || "Peringatan Stok");
      const message = escapeHTML(notification.message || "-");
      card.innerHTML = `
        <div class="flex items-start gap-3">
          <div class="mt-1 flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-orange-50 text-orange-600">
            <i data-lucide="alert-triangle" class="h-5 w-5"></i>
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-black text-slate-900">${title}</p>
            <p class="mt-1 text-sm leading-5 text-slate-600">${message}</p>
            <div class="mt-3 flex flex-wrap gap-2">
              <button type="button" data-action="open-products" class="rounded-lg bg-slate-900 px-3 py-2 text-xs font-bold text-white hover:bg-slate-800">
                Buka Produk Saya
              </button>
              <button type="button" data-action="read" class="rounded-lg border border-slate-200 px-3 py-2 text-xs font-bold text-slate-600 hover:bg-slate-50">
                Tandai Dibaca
              </button>
            </div>
          </div>
        </div>
      `;

      card.querySelector('[data-action="open-products"]').addEventListener("click", () => {
        const query = notification.source_id
          ? `?q=${encodeURIComponent(notification.source_id)}`
          : "";
        window.location.href = `supplier_produk_saya.html${query}`;
      });
      card.querySelector('[data-action="read"]').addEventListener("click", async () => {
        try {
          await markNotificationRead(notification.id);
          card.remove();
          if (!panel.children.length) panel.remove();
        } catch (error) {
          notify("danger", error.message || "Gagal menandai notifikasi");
        }
      });

      panel.appendChild(card);
    });

    document.body.appendChild(panel);
    if (window.lucide) lucide.createIcons();
  }

  async function showUnreadNotifications() {
    try {
      const notifications = await getUnreadNotifications();
      renderNotificationPanel(notifications);
      if (notifications.length > 0) {
        notify("warning", `${notifications.length} peringatan stok belum dibaca.`);
      }
    } catch (error) {
      console.warn("Gagal mengambil notifikasi supplier", error);
    }
  }

  async function initPage(options = {}) {
    const user = requireSupplier(options.redirectUrl);
    if (!user) return null;

    setUserDisplay(user);
    setCurrentDate();
    setActiveNavigation();

    const appContent = document.getElementById("app-content");
    if (appContent) appContent.classList.remove("invisible");
    if (window.lucide) lucide.createIcons();

    if (typeof options.onReady === "function") {
      await options.onReady(user);
    }

    if (!options.skipUnreadNotifications) {
      await showUnreadNotifications();
    }

    return user;
  }

  window.toggleSidebar = toggleSidebar;

  return {
    apiFetch,
    escapeHTML,
    formatNumber,
    formatRupiah,
    getOrders,
    getProducts,
    getProfile,
    getStats,
    getNotifications,
    getUnreadNotifications,
    initPage,
    markNotificationRead,
    renderNotificationPanel,
    notify,
    setUserDisplay,
    shortID,
    statusMeta,
    toggleSidebar,
    updateOrderStatus,
    updateProfile,
    updateStoredProfile,
  };
})();
