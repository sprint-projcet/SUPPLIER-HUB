// auth.js
// Simulasi Autentikasi dan Konfigurasi Dashboard Multi-Role

const roleConfig = {
    admin: {
        title: 'Admin Command Center',
        stats: [
            { label: 'Total Supplier', value: '128', icon: 'users', change: '+12%', positive: true },
            { label: 'Total Transaksi', value: 'Rp 2.4M', icon: 'shopping-cart', change: '+18%', positive: true },
            { label: 'Pesanan Aktif', value: '45', icon: 'package', change: '-4%', positive: false },
            { label: 'Revenue Growth', value: '24%', icon: 'bar-chart-3', change: '+2%', positive: true },
        ],
        nav: [
            { name: 'Overview', icon: 'layout-dashboard' },
            { name: 'Daftar Supplier', icon: 'users' },
            { name: 'Kontrol Stok', icon: 'package' },
            { name: 'Keuangan', icon: 'bar-chart-3' },
            { name: 'Pengaturan', icon: 'settings' }
        ],
        tableTitle: 'Monitoring Transaksi Global',
        quickBtn: 'Supplier'
    },
    supplier: {
        title: 'Supplier Portal',
        stats: [
            { label: 'Stok Barang', value: '1,240', icon: 'package', change: '+5%', positive: true },
            { label: 'Pesanan Baru', value: '12', icon: 'bell', change: 'Baru', positive: true },
            { label: 'Pendapatan', value: 'Rp 450jt', icon: 'bar-chart-3', change: '+10%', positive: true },
            { label: 'Rating Toko', value: '4.8/5', icon: 'check-circle-2', change: 'Stabil', positive: true },
        ],
        nav: [
            { name: 'Dashboard', icon: 'layout-dashboard' },
            { name: 'Produk Saya', icon: 'package' },
            { name: 'Daftar Pesanan', icon: 'shopping-cart' },
            { name: 'Analitik Toko', icon: 'bar-chart-3' },
            { name: 'Toko Saya', icon: 'settings' }
        ],
        tableTitle: 'Pesanan Masuk Terbaru',
        quickBtn: 'Produk'
    },
    user: {
        title: 'Dashboard UMKM (Pembeli)',
        stats: [
            { label: 'Total Pesanan', value: '24', icon: 'shopping-cart', change: 'Bulan ini', positive: true },
            { label: 'Sedang Dikirim', value: '3', icon: 'truck', change: 'Aktif', positive: true },
            { label: 'Voucher', value: '5', icon: 'tag', change: 'Tersedia', positive: true },
            { label: 'Poin Hub', value: '12,500', icon: 'zap', change: '+500', positive: true },
        ],
        nav: [
            { name: 'Belanja', icon: 'layout-dashboard' },
            { name: 'Pesanan Saya', icon: 'shopping-cart' },
            { name: 'Lacak Paket', icon: 'truck' },
            { name: 'Wishlist', icon: 'heart' },
            { name: 'Bantuan', icon: 'help-circle' }
        ],
        tableTitle: 'Riwayat Pembelian Saya',
        quickBtn: 'Order'
    }
};

/**
 * Fungsi untuk simulasi login ke sistem.
 */
function loginUser(email, password, role) {
    return fetch("http://localhost:8080/api/auth/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ email: email, password: password })
    })
    .then(response => response.json().then(data => ({ status: response.status, ok: response.ok, data })))
    .then(({ status, ok, data }) => {
        if (!ok) {
            throw new Error(data.error || 'Terjadi kesalahan saat login.');
        }
        
        const userSession = {
            name: data.user.business_name || role,
            role: data.role,
            email: data.user.email,
            token: data.token,
            id: data.user.id,
            lastLogin: new Date().toISOString()
        };
        
        localStorage.setItem('user_session', JSON.stringify(userSession));
        return userSession;
    });
}

/**
 * Memastikan sesi pengguna ada. Jika tidak, redirect ke halaman login.
 */
function checkAuth(redirectUrl = null) {
    const sessionData = localStorage.getItem('user_session');
    if (!sessionData || sessionData === "undefined" || sessionData === "null") {
        if (redirectUrl) window.location.href = redirectUrl;
        return null;
    }
    
    try {
        const user = JSON.parse(sessionData);
        return user;
    } catch (e) {
        console.error("Data sesi korup:", e);
        localStorage.removeItem('user_session');
        if (redirectUrl) window.location.href = redirectUrl;
        return null;
    }
}

/**
 * Fungsi Logout untuk memutus sesi.
 */
function logoutUser(redirectUrl = '../Login/login.html') {
    localStorage.removeItem('user_session');
    window.location.href = redirectUrl;
}

/**
 * Mengambil konfigurasi antarmuka/data untuk tiap role dashboard
 */
function getRoleConfig(role) {
    return roleConfig[role] || roleConfig['user'];
}
