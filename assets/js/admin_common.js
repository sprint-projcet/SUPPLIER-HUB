const AdminDashboard = (() => {
  function getSession() {
    return typeof getStoredUserSession === "function"
      ? getStoredUserSession()
      : null;
  }

  function requireAdmin(redirectUrl = "../Login/login.html") {
    const user = checkAuth(redirectUrl);
    if (!user) return null;

    if (user.role !== "admin") {
      window.location.href = redirectUrl;
      return null;
    }

    return user;
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
      const error = new Error(data.error || "Request admin gagal diproses");
      error.status = response.status;
      error.data = data;
      throw error;
    }

    return data;
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

  function formatDate(value) {
    if (!value) return "-";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "-";
    return date.toLocaleDateString("id-ID", {
      day: "numeric",
      month: "short",
      year: "numeric",
    });
  }

  function escapeHTML(value) {
    return String(value ?? "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  }

  function shortID(prefix, value) {
    const id = String(value || "");
    if (!id) return `${prefix}-`;
    return `${prefix}-${id.slice(0, 8).toUpperCase()}`;
  }

  function orderStatusMeta(status) {
    const normalized = String(status || "").toLowerCase();
    const map = {
      pending: ["PENDING", "bg-yellow-50 text-yellow-600"],
      pending_supplier_confirmation: [
        "MENUNGGU SUPPLIER",
        "bg-orange-50 text-orange-600",
      ],
      rejected_by_supplier: ["DITOLAK", "bg-red-50 text-red-600"],
      stock_unavailable: ["STOK HABIS", "bg-red-50 text-red-600"],
      supplier_confirmed: ["STOK OK", "bg-teal-50 text-teal-600"],
      payment_pending: ["MENUNGGU BAYAR", "bg-amber-50 text-amber-600"],
      payment_request_failed: ["PAYMENT GAGAL", "bg-red-50 text-red-600"],
      paid: ["LUNAS", "bg-emerald-50 text-emerald-600"],
      payment_failed: ["GAGAL BAYAR", "bg-red-50 text-red-600"],
      shipment_created: ["PENGIRIMAN", "bg-indigo-50 text-indigo-600"],
      processing: ["PROSES", "bg-blue-50 text-blue-600"],
      shipped: ["DIKIRIM", "bg-indigo-50 text-indigo-600"],
      completed: ["SELESAI", "bg-emerald-50 text-emerald-600"],
      cancelled: ["BATAL", "bg-red-50 text-red-600"],
    };

    return map[normalized] || [
      normalized ? normalized.toUpperCase() : "-",
      "bg-slate-100 text-slate-600",
    ];
  }

  function supplierStatusMeta(status) {
    const normalized = String(status || "").toLowerCase();
    const map = {
      active: ["AKTIF", "bg-emerald-50 text-emerald-600"],
      pending: ["MENUNGGU VERIFIKASI", "bg-slate-100 text-slate-600"],
      suspended: ["DITANGGUHKAN", "bg-red-50 text-red-600"],
    };

    return map[normalized] || [
      normalized ? normalized.toUpperCase() : "-",
      "bg-slate-100 text-slate-600",
    ];
  }

  function badge(label, className) {
    return `<span class="text-[10px] font-black px-3 py-1 rounded-full ${className}">${escapeHTML(label)}</span>`;
  }

  function notify(type, message) {
    if (typeof window.showGlobalToast === "function") {
      window.showGlobalToast(type, message);
      return;
    }
    if (typeof window.showToast === "function") {
      window.showToast(type, message);
      return;
    }
    alert(message);
  }

  function getNotifications(options = {}) {
    const query = options.unreadOnly ? "?unread_only=true" : "";
    return apiFetch(`/api/admin/notifications${query}`).then(
      (result) => result.data || [],
    );
  }

  function markNotificationRead(notificationID) {
    return apiFetch(`/api/admin/notifications/${encodeURIComponent(notificationID)}/read`, {
      method: "PUT",
      body: JSON.stringify({}),
    });
  }

  function renderNotificationPanel(notifications) {
    document.getElementById("admin-notification-panel")?.remove();
    if (!Array.isArray(notifications) || notifications.length === 0) return;

    const panel = document.createElement("div");
    panel.id = "admin-notification-panel";
    panel.className =
      "fixed left-1/2 transform -translate-x-1/2 top-24 z-[120] w-[calc(100%-2rem)] max-w-sm space-y-2";

    notifications.slice(0, 5).forEach((notification, index) => {
      const card = document.createElement("div");
      card.className =
        "notification-toast flex items-center gap-3 px-4 py-3 rounded-lg bg-white border border-slate-200 shadow-md backdrop-blur-sm animate-slide-in";
      card.style.animationDelay = `${index * 100}ms`;

      const title = escapeHTML(notification.title || "Notifikasi Admin");
      const message = escapeHTML(notification.message || "-");
      card.innerHTML = `
        <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600">
          <i data-lucide="check-circle" class="h-4 w-4"></i>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-semibold text-slate-900">${title}</p>
          <p class="text-xs text-slate-500 leading-relaxed line-clamp-2">${message}</p>
        </div>
        <button type="button" data-action="close" class="shrink-0 text-slate-400 hover:text-slate-600 transition-colors">
          <i data-lucide="x" class="h-4 w-4"></i>
        </button>
      `;

      card.querySelector('[data-action="close"]').addEventListener("click", async () => {
        card.classList.add("animate-slide-out");
        setTimeout(() => {
          card.remove();
          if (!panel.children.length) panel.remove();
        }, 300);
        try {
          await markNotificationRead(notification.id);
        } catch (error) {
          console.error("Gagal menandai notifikasi", error);
        }
      });

      card.addEventListener("click", (e) => {
        if (e.target.closest('[data-action="close"]')) return;
        const query = notification.source_id
          ? `?q=${encodeURIComponent(notification.source_id)}`
          : "";
        window.location.href = `admin_kontrol_stok.html${query}`;
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
      const notifications = await getNotifications({ unreadOnly: true });
      renderNotificationPanel(notifications);
      if (notifications.length > 0) {
        notify("info", `${notifications.length} notifikasi admin belum dibaca.`);
      }
    } catch (error) {
      console.warn("Gagal mengambil notifikasi admin", error);
    }
  }

  function downloadCSV(filename, header, rows) {
    const csvRows = [
      header,
      ...rows.map((row) =>
        row
          .map((value) => `"${String(value ?? "").replace(/"/g, '""')}"`)
          .join(","),
      ),
    ];

    const blob = new Blob([csvRows.join("\n")], {
      type: "text/csv;charset=utf-8;",
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = filename;
    link.click();
    URL.revokeObjectURL(url);
  }

  async function initPage(options = {}) {
    const user = requireAdmin(options.redirectUrl);
    if (!user) return null;

    const displayName = user.business_name || user.name || user.email || "Admin";
    document.querySelectorAll("#user-display-name").forEach((element) => {
      element.textContent = displayName;
    });

    const dateEl = document.getElementById("current-date");
    if (dateEl) {
      dateEl.textContent = new Date().toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
      });
    }

    const content = document.getElementById("app-content");
    if (content) content.classList.remove("invisible");

    if (window.lucide) lucide.createIcons();

    if (typeof options.onReady === "function") {
      await options.onReady(user);
    }

    await showUnreadNotifications();

    return user;
  }

  window.toggleSidebar = function toggleSidebar() {
    const sidebar = document.getElementById("sidebar");
    if (sidebar) sidebar.classList.toggle("-translate-x-full");
  };

  return {
    apiFetch,
    badge,
    downloadCSV,
    escapeHTML,
    formatDate,
    formatNumber,
    formatRupiah,
    getNotifications,
    initPage,
    markNotificationRead,
    notify,
    orderStatusMeta,
    renderNotificationPanel,
    showUnreadNotifications,
    shortID,
    supplierStatusMeta,
  };
})();
