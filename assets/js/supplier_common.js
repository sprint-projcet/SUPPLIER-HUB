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
    ensureSupplierChatNavigation();
    const currentPage = window.location.pathname.split("/").pop() || "supplier.html";

    document.querySelectorAll("#sidebar nav a[href]").forEach((link) => {
      const isActive = link.getAttribute("href") === currentPage;
      link.className = `${NAV_BASE_CLASS} ${isActive ? NAV_ACTIVE_CLASS : NAV_IDLE_CLASS}`;

      link.querySelectorAll("i").forEach((icon) => {
        icon.className = isActive ? "text-white" : "group-hover:text-emerald-400";
      });
    });
  }

  function ensureSupplierChatNavigation() {
    const nav = document.querySelector("#sidebar nav");
    if (!nav || nav.querySelector('a[href="supplier_chat.html"]')) return;

    const link = document.createElement("a");
    link.href = "supplier_chat.html";
    link.className = `${NAV_BASE_CLASS} ${NAV_IDLE_CLASS}`;
    link.innerHTML = `
      <i data-lucide="message-circle" class="group-hover:text-emerald-400"></i>
      <span class="font-semibold text-sm">Chat UMKM</span>
    `;

    const beforeLink =
      nav.querySelector('a[href="supplier_notifikasi.html"]') ||
      nav.querySelector('a[href="supplier_analitik.html"]') ||
      nav.querySelector('a[href="supplier_toko.html"]');
    nav.insertBefore(link, beforeLink || null);
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

  function formatNotificationTime(value) {
    if (!value) return "-";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "-";
    return date.toLocaleString("id-ID", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
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
      "fixed left-1/2 transform -translate-x-1/2 top-24 z-[120] w-[calc(100%-2rem)] max-w-sm space-y-2";

    notifications.slice(0, 5).forEach((notification, index) => {
      const isChatNotification =
        notification.source_type === "chat" || notification.type === "chat_message";
      const card = document.createElement("div");
      card.className =
        "notification-toast flex items-center gap-3 px-4 py-3 rounded-lg bg-white border border-slate-200 shadow-md backdrop-blur-sm animate-slide-in cursor-pointer";
      card.style.animationDelay = `${index * 100}ms`;

      const title = escapeHTML(notification.title || "Peringatan Stok");
      const message = escapeHTML(notification.message || "-");
      const createdAt = escapeHTML(formatNotificationTime(notification.created_at));
      card.innerHTML = `
        <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg ${isChatNotification ? "bg-emerald-50 text-emerald-600" : "bg-orange-50 text-orange-600"}">
          <i data-lucide="${isChatNotification ? "message-circle" : "alert-circle"}" class="h-4 w-4"></i>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-semibold text-slate-900">${title}</p>
          <p class="text-xs text-slate-500 leading-relaxed line-clamp-2">${message}</p>
          <p class="mt-1 text-[10px] font-bold uppercase tracking-widest text-slate-400">${createdAt}</p>
        </div>
        <button type="button" data-action="close" class="shrink-0 text-slate-400 hover:text-slate-600 transition-colors">
          <i data-lucide="x" class="h-4 w-4"></i>
        </button>
      `;

      card.querySelector('[data-action="close"]').addEventListener("click", (e) => {
        e.stopPropagation();
        card.classList.add("animate-slide-out");
        setTimeout(() => {
          card.remove();
          if (!panel.children.length) panel.remove();
        }, 300);
        try {
          markNotificationRead(notification.id);
        } catch (error) {
          console.error("Gagal menandai notifikasi", error);
        }
      });

      card.addEventListener("click", (e) => {
        if (e.target.closest('[data-action="close"]')) return;
        if (notification.source_type === "chat") {
          const query = notification.source_id
            ? `?conversation_id=${encodeURIComponent(notification.source_id)}`
            : "";
          window.location.href = `supplier_chat.html${query}`;
          return;
        }
        const query = notification.source_id
          ? `?q=${encodeURIComponent(notification.source_id)}`
          : "";
        window.location.href = `supplier_produk_saya.html${query}`;
      });

      const autoHideTimer = setTimeout(() => {
        card.classList.add("animate-slide-out");
        setTimeout(() => {
          card.remove();
          if (!panel.children.length) panel.remove();
        }, 300);
      }, 5000);

      card.addEventListener("mouseenter", () => clearTimeout(autoHideTimer));

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
