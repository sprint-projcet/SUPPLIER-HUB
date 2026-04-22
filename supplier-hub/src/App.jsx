import React, { useState, useMemo } from 'react';
import { 
  LayoutDashboard, Package, ShoppingCart, Wallet, Truck, 
  Settings, Search, Plus, ChevronRight, CheckCircle2, 
  Menu, X, ArrowUpRight, Star, LogOut, User, 
  Users, BarChart3, ShieldCheck, ShoppingBag, Factory, Info
} from 'lucide-react';

// --- KONFIGURASI WARNA ---
const COLORS = {
  primary: '#028090', // Teal Utama
  secondary: '#02C39A', // Teal Muda
  navy: '#0D1B2A',     // Navy
  bg: '#F8FAFC'
};

// --- DATA PRODUK (MOCK) ---
const MOCK_PRODUCTS = [
  { id: 1, name: 'Kain Katun Premium', price: 45000, category: 'Tekstil', unit: 'Meter' },
  { id: 2, name: 'Benang Sutra', price: 12000, category: 'Tekstil', unit: 'Roll' },
  { id: 3, name: 'Kabel Tembaga 1m', price: 8500, category: 'Elektronik', unit: 'Pcs' },
  { id: 4, name: 'Kardus Box Custom', price: 3500, category: 'Kemasan', unit: 'Pcs' },
];

// --- KOMPONEN LOGO ---
const Logo = () => (
  <div className="flex items-center gap-3">
    <div className="relative w-10 h-10 flex items-center justify-center">
      <svg viewBox="0 0 100 100" className="absolute w-full h-full drop-shadow-md">
        <polygon points="50,5 95,25 95,75 50,95 5,75 5,25" fill={COLORS.primary} />
        <rect x="35" y="35" width="30" height="30" transform="rotate(45 50 50)" fill="white" fillOpacity="0.3" />
      </svg>
    </div>
    <div className="flex flex-col">
      <span className="wordmark text-xl font-bold leading-none">
        <span style={{ color: COLORS.navy }}>Supplier</span>
        <span style={{ color: COLORS.primary }}>Hub</span>
      </span>
      <span className="tagline text-[10px] text-slate-400 uppercase tracking-tighter">B2B Ecosystem</span>
    </div>
  </div>
);

export default function App() {
  // --- STATE UTAMA ---
  const [role, setRole] = useState('UMKM'); // Switcher: 'UMKM', 'SUPPLIER', 'ADMIN'
  const [activePage, setActivePage] = useState('dashboard');
  const [isSidebarOpen, setSidebarOpen] = useState(true);
  const [cart, setCart] = useState([]);
  const [toast, setToast] = useState({ show: false, msg: '' });

  // Data User Sesuai Role
  const userData = {
    UMKM: { name: "Toko Berkah UMKM", label: "Gold Member", balance: 24500000 },
    SUPPLIER: { name: "PT. Maju Tekstil", label: "Verified Supplier", sales: "Rp 120jt" },
    ADMIN: { name: "Central Admin", label: "Superuser", alerts: 0 }
  };

  const current = userData[role];

  // Logika Keranjang (Hanya UMKM)
  const cartSubtotal = useMemo(() => cart.reduce((acc, item) => acc + (item.price * item.qty), 0), [cart]);
  const feeAdmin = cartSubtotal * 0.03;
  const totalPay = cartSubtotal + feeAdmin;

  const addToCart = (product) => {
    setCart(prev => {
      const exist = prev.find(p => p.id === product.id);
      if (exist) return prev.map(p => p.id === product.id ? { ...p, qty: p.qty + 1 } : p);
      return [...prev, { ...product, qty: 1 }];
    });
    setToast({ show: true, msg: `${product.name} masuk keranjang!` });
    setTimeout(() => setToast({ show: false, msg: '' }), 2000);
  };

  return (
    <div className="min-h-screen flex bg-[#F8FAFC] text-slate-900 font-sans">
      
      {/* SIDEBAR */}
      <aside className={`fixed inset-y-0 left-0 z-50 w-72 bg-white border-r border-slate-100 transition-transform lg:relative lg:translate-x-0 ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div className="p-8 flex flex-col h-full">
          <Logo />
          <nav className="mt-12 flex-1 space-y-2">
            <button onClick={() => setActivePage('dashboard')} className={`w-full flex items-center gap-4 px-4 py-3 rounded-2xl font-bold transition-all ${activePage === 'dashboard' ? 'bg-teal-50 text-[#028090]' : 'text-slate-400 hover:bg-slate-50'}`}>
              <LayoutDashboard size={20} /> Dashboard
            </button>

            {role === 'UMKM' && (
              <button onClick={() => setActivePage('catalog')} className={`w-full flex items-center gap-4 px-4 py-3 rounded-2xl font-bold transition-all ${activePage === 'catalog' ? 'bg-teal-50 text-[#028090]' : 'text-slate-400 hover:bg-slate-50'}`}>
                <Search size={20} /> Cari Bahan
              </button>
            )}

            {role === 'SUPPLIER' && (
              <button onClick={() => setActivePage('inventory')} className={`w-full flex items-center gap-4 px-4 py-3 rounded-2xl font-bold transition-all ${activePage === 'inventory' ? 'bg-teal-50 text-[#028090]' : 'text-slate-400 hover:bg-slate-50'}`}>
                <Factory size={20} /> Kelola Produk
              </button>
            )}

            {role === 'ADMIN' && (
              <button onClick={() => setActivePage('users')} className={`w-full flex items-center gap-4 px-4 py-3 rounded-2xl font-bold transition-all ${activePage === 'users' ? 'bg-teal-50 text-[#028090]' : 'text-slate-400 hover:bg-slate-50'}`}>
                <Users size={20} /> Data Member
              </button>
            )}

            <button className="w-full flex items-center gap-4 px-4 py-3 rounded-2xl font-bold text-slate-400 hover:bg-slate-50">
              <Wallet size={20} /> SmartBank
            </button>
          </nav>
          <div className="pt-6 border-t border-slate-100 flex items-center gap-3">
            <div className="w-10 h-10 bg-[#0D1B2A] text-white rounded-xl flex items-center justify-center font-bold text-xs uppercase">{role[0]}</div>
            <div className="flex-1 overflow-hidden">
              <p className="text-sm font-bold truncate">{current.name}</p>
              <p className="text-[10px] text-teal-600 font-bold uppercase">{role}</p>
            </div>
          </div>
        </div>
      </aside>

      {/* MAIN CONTENT */}
      <main className="flex-1 h-screen overflow-hidden flex flex-col">
        {/* Header */}
        <header className="h-20 bg-white border-b border-slate-100 px-8 flex items-center justify-between shadow-sm">
          <div className="flex items-center gap-4">
            <button onClick={() => setSidebarOpen(!isSidebarOpen)} className="lg:hidden p-2"><Menu /></button>
            <div className="bg-slate-100 p-1 rounded-xl flex gap-1">
              {['UMKM', 'SUPPLIER', 'ADMIN'].map(r => (
                <button key={r} onClick={() => {setRole(r); setActivePage('dashboard');}} className={`px-3 py-1 rounded-lg text-[10px] font-black transition-all ${role === r ? 'bg-white text-teal-600 shadow-sm' : 'text-slate-400'}`}>{r}</button>
              ))}
            </div>
          </div>
          <div className="flex items-center gap-6">
            {role === 'UMKM' && (
              <div className="relative cursor-pointer" onClick={() => setActivePage('cart')}>
                <ShoppingCart size={22} className="text-slate-600" />
                {cart.length > 0 && <span className="absolute -top-2 -right-2 bg-orange-500 text-white text-[10px] w-5 h-5 flex items-center justify-center rounded-full border-2 border-white">{cart.length}</span>}
              </div>
            )}
            <div className="w-8 h-8 bg-teal-100 text-teal-700 rounded-lg flex items-center justify-center font-bold"><User size={18} /></div>
          </div>
        </header>

        {/* Dashboard/Content */}
        <section className="flex-1 overflow-y-auto p-8 custom-scrollbar">
          <div className="max-w-6xl mx-auto">
            <h1 className="text-3xl font-black mb-8 tracking-tight">Dashboard {role}</h1>
            
            {activePage === 'dashboard' && (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm">
                  <p className="text-slate-500 text-sm mb-1">Status Keanggotaan</p>
                  <p className="text-2xl font-black text-teal-600">{current.label}</p>
                </div>
                <div className="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm">
                  <p className="text-slate-500 text-sm mb-1">{role === 'UMKM' ? 'Saldo SmartBank' : 'Total Penjualan'}</p>
                  <p className="text-2xl font-black text-slate-900">{role === 'UMKM' ? `Rp ${current.balance.toLocaleString()}` : current.sales}</p>
                </div>
                <div className="bg-[#0D1B2A] text-white p-6 rounded-3xl relative overflow-hidden">
                  <p className="text-teal-400 text-[10px] font-bold uppercase mb-1">Notifikasi Sistem</p>
                  <p className="text-2xl font-black">0 Pesan Baru</p>
                  <div className="absolute -bottom-5 -right-5 w-20 h-20 bg-teal-500/20 rounded-full blur-2xl"></div>
                </div>
              </div>
            )}

            {activePage === 'catalog' && (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 animate-in slide-in-from-bottom-4">
                {MOCK_PRODUCTS.map(p => (
                  <div key={p.id} className="bg-white p-6 rounded-3xl border border-slate-100 hover:shadow-xl transition-all">
                    <div className="aspect-square bg-slate-50 rounded-2xl mb-4 flex items-center justify-center text-slate-200"><Package size={48} /></div>
                    <p className="text-[10px] font-bold text-teal-600 uppercase mb-1">{p.category}</p>
                    <h4 className="font-bold text-slate-900 mb-4">{p.name}</h4>
                    <div className="flex justify-between items-center pt-4 border-t border-slate-50">
                      <p className="font-black text-sm">Rp {p.price.toLocaleString()}</p>
                      <button onClick={() => addToCart(p)} className="p-2 bg-[#0D1B2A] text-white rounded-xl hover:bg-teal-600 transition-colors"><Plus size={18}/></button>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {activePage === 'cart' && (
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 space-y-4">
                  {cart.map(item => (
                    <div key={item.id} className="bg-white p-4 rounded-2xl border border-slate-100 flex items-center gap-4">
                      <div className="flex-1 font-bold">{item.name} (x{item.qty})</div>
                      <p className="font-black">Rp {(item.qty * item.price).toLocaleString()}</p>
                      <button onClick={() => setCart(c => c.filter(i => i.id !== item.id))} className="text-slate-300 hover:text-red-500"><X size={18}/></button>
                    </div>
                  ))}
                  {cart.length === 0 && <p className="text-center py-20 text-slate-400 italic">Keranjang kosong</p>}
                </div>
                <div className="bg-white p-8 rounded-3xl border border-slate-100 shadow-xl h-fit">
                  <h3 className="font-bold mb-6">Total Pembayaran</h3>
                  <div className="flex justify-between mb-4 text-slate-500 text-sm"><span>Subtotal</span><span>Rp {cartSubtotal.toLocaleString()}</span></div>
                  <div className="flex justify-between mb-8 text-slate-500 text-sm"><span>Layanan (3%)</span><span>Rp {feeAdmin.toLocaleString()}</span></div>
                  <div className="flex justify-between items-end mb-8"><span className="font-bold">Total</span><span className="text-2xl font-black text-teal-600">Rp {totalPay.toLocaleString()}</span></div>
                  <button className="w-full py-4 bg-teal-600 text-white rounded-2xl font-black shadow-lg hover:bg-[#0D1B2A] transition-all disabled:bg-slate-200" disabled={cart.length === 0}>Bayar Sekarang</button>
                </div>
              </div>
            )}
          </div>
        </section>
      </main>

      {/* TOAST */}
      {toast.show && (
        <div className="fixed bottom-10 left-1/2 -translate-x-1/2 bg-[#0D1B2A] text-white px-6 py-3 rounded-2xl shadow-2xl flex items-center gap-3 animate-in slide-in-from-bottom-4">
          <CheckCircle2 size={18} className="text-teal-400" />
          <span className="font-bold text-sm">{toast.msg}</span>
        </div>
      )}
    </div>
  );
}