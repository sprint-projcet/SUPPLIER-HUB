// auth.js
// Autentikasi dan Konfigurasi Dashboard Multi-Role

const API_BASE_URL =
  window.SUPPLIER_HUB_API_BASE_URL || "http://localhost:8080";
const LEGACY_AUTH_KEYS = ["authToken", "userRole", "userData", "token"];

function buildApiUrl(path) {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  return `${API_BASE_URL}${normalizedPath}`;
}

const roleConfig = {
  admin: {
    title: "Admin Command Center",
    stats: [
      {
        label: "Total Supplier",
        value: "128",
        icon: "users",
        change: "+12%",
        positive: true,
      },
      {
        label: "Total Transaksi",
        value: "Rp 2.4M",
        icon: "shopping-cart",
        change: "+18%",
        positive: true,
      },
      {
        label: "Pesanan Aktif",
        value: "45",
        icon: "package",
        change: "-4%",
        positive: false,
      },
      {
        label: "Revenue Growth",
        value: "24%",
        icon: "bar-chart-3",
        change: "+2%",
        positive: true,
      },
    ],
    nav: [
      { name: "Overview", icon: "layout-dashboard" },
      { name: "Daftar Supplier", icon: "users" },
      { name: "Kontrol Stok", icon: "package" },
      { name: "Keuangan", icon: "bar-chart-3" },
      { name: "Pengaturan", icon: "settings" },
    ],
    tableTitle: "Monitoring Transaksi Global",
    quickBtn: "Supplier",
  },
  supplier: {
    title: "Supplier Portal",
    stats: [
      {
        label: "Stok Barang",
        value: "1,240",
        icon: "package",
        change: "+5%",
        positive: true,
      },
      {
        label: "Pesanan Baru",
        value: "12",
        icon: "bell",
        change: "Baru",
        positive: true,
      },
      {
        label: "Pendapatan",
        value: "Rp 450jt",
        icon: "bar-chart-3",
        change: "+10%",
        positive: true,
      },
      {
        label: "Rating Toko",
        value: "4.8/5",
        icon: "check-circle-2",
        change: "Stabil",
        positive: true,
      },
    ],
    nav: [
      { name: "Dashboard", icon: "layout-dashboard" },
      { name: "Produk Saya", icon: "package" },
      { name: "Daftar Pesanan", icon: "shopping-cart" },
      { name: "Analitik Toko", icon: "bar-chart-3" },
      { name: "Toko Saya", icon: "settings" },
    ],
    tableTitle: "Pesanan Masuk Terbaru",
    quickBtn: "Produk",
  },
  user: {
    title: "Dashboard UMKM (Pembeli)",
    stats: [
      {
        label: "Total Pesanan",
        value: "24",
        icon: "shopping-cart",
        change: "Bulan ini",
        positive: true,
      },
      {
        label: "Sedang Dikirim",
        value: "3",
        icon: "truck",
        change: "Aktif",
        positive: true,
      },
      {
        label: "Voucher",
        value: "5",
        icon: "tag",
        change: "Tersedia",
        positive: true,
      },
      {
        label: "Poin Hub",
        value: "12,500",
        icon: "zap",
        change: "+500",
        positive: true,
      },
    ],
    nav: [
      { name: "Belanja", icon: "layout-dashboard" },
      { name: "Pesanan Saya", icon: "shopping-cart" },
      { name: "Lacak Paket", icon: "truck" },
      { name: "Wishlist", icon: "heart" },
      { name: "Bantuan", icon: "help-circle" },
    ],
    tableTitle: "Riwayat Pembelian Saya",
    quickBtn: "Order",
  },
};

/**
 * Menormalkan response auth backend menjadi format sesi frontend.
 */
function createUserSession(data, fallbackRole = "user") {
  const user = data.user || {};
  const role = data.role || fallbackRole;

  return {
    name: user.business_name || user.email || role,
    role: role,
    email: user.email || "",
    token: data.token,
    id: user.id || "",
    lastLogin: new Date().toISOString(),
  };
}

/**
 * Menyimpan sesi di satu sumber data agar dashboard membaca format yang konsisten.
 */
function saveUserSession(userSession) {
  localStorage.setItem("user_session", JSON.stringify(userSession));
  LEGACY_AUTH_KEYS.forEach((key) => localStorage.removeItem(key));
  return userSession;
}

function saveUserSessionFromAuthResponse(data, fallbackRole = "user") {
  return saveUserSession(createUserSession(data, fallbackRole));
}

function clearUserSession() {
  localStorage.removeItem("user_session");
  LEGACY_AUTH_KEYS.forEach((key) => localStorage.removeItem(key));
}

function getStoredUserSession() {
  const sessionData = localStorage.getItem("user_session");
  if (!sessionData || sessionData === "undefined" || sessionData === "null") {
    return null;
  }

  try {
    return JSON.parse(sessionData);
  } catch (e) {
    console.error("Data sesi korup:", e);
    clearUserSession();
    return null;
  }
}

/**
 * Fungsi untuk login ke sistem.
 */
function loginUser(email, password, role) {
  return fetch(buildApiUrl("/api/auth/login"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email: email, password: password }),
  })
    .then((response) =>
      response
        .json()
        .then((data) => ({ status: response.status, ok: response.ok, data })),
    )
    .then(({ status, ok, data }) => {
      if (!ok) {
        throw new Error(data.error || "Terjadi kesalahan saat login.");
      }

      return saveUserSessionFromAuthResponse(data, role);
    });
}

/**
 * Memastikan sesi pengguna ada. Jika tidak, redirect ke halaman login.
 */
function checkAuth(redirectUrl = null) {
  const user = getStoredUserSession();
  if (!user) {
    if (redirectUrl) window.location.href = redirectUrl;
    return null;
  }

  return user;
}

/**
 * Fungsi Logout untuk memutus sesi.
 */
function logoutUser(redirectUrl = "../Login/login.html") {
  clearUserSession();
  sessionStorage.setItem("justLoggedOut", "true");
  window.location.href = redirectUrl;
}

/**
 * Mengambil konfigurasi antarmuka/data untuk tiap role dashboard
 */
function getRoleConfig(role) {
  return roleConfig[role] || roleConfig["user"];
}
