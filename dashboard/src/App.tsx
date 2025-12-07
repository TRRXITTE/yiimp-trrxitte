import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import axios from 'axios';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const API_URL = 'http://localhost:8001/api';

interface PoolStats {
  pool_name: string;
  hashrate: number;
  miners: number;
  workers: number;
  blocks_found: number;
  last_block_time: string;
  network_difficulty: { [key: string]: number };
  pool_difficulty: number;
  sync_status: SyncStatus[];
  chains: string[];
  merged_mining: boolean;
}

interface SyncStatus {
  chain: string;
  blocks: number;
  headers: number;
  sync_progress: number;
  is_syncing: boolean;
  verification_progress: number;
}

interface MinerStats {
  address: string;
  hashrate: number;
  shares_valid: number;
  shares_invalid: number;
  last_share: string;
  balance: { [key: string]: number };
  paid: { [key: string]: number };
  workers: WorkerStats[];
}

interface WorkerStats {
  name: string;
  hashrate: number;
  shares_valid: number;
  shares_invalid: number;
  last_share: string;
}

interface Block {
  height: number;
  hash: string;
  chain: string;
  difficulty: number;
  reward: number;
  confirmations: number;
  status: string;
  time: string;
  miner: string;
}

function Dashboard() {
  const [stats, setStats] = useState<PoolStats | null>(null);
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [statsRes, blocksRes] = await Promise.all([
          axios.get(`${API_URL}/pool/stats`),
          axios.get(`${API_URL}/pool/blocks?limit=10`)
        ]);
        setStats(statsRes.data);
        setBlocks(blocksRes.data);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching data:', error);
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 10000); // Update every 10 seconds

    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-blue-600 text-white shadow-lg">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">{stats?.pool_name} Mining Pool</h1>
            <div className="flex space-x-6">
              <Link to="/" className="hover:text-blue-200">Dashboard</Link>
              <Link to="/miners" className="hover:text-blue-200">Miners</Link>
              <Link to="/blocks" className="hover:text-blue-200">Blocks</Link>
              <Link to="/payments" className="hover:text-blue-200">Payments</Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="container mx-auto px-6 py-8">
        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatCard
            title="Pool Hashrate"
            value={formatHashrate(stats?.hashrate || 0)}
            icon="⚡"
            color="blue"
          />
          <StatCard
            title="Active Miners"
            value={stats?.miners.toString() || '0'}
            icon="👥"
            color="green"
          />
          <StatCard
            title="Active Workers"
            value={stats?.workers.toString() || '0'}
            icon="🔧"
            color="purple"
          />
          <StatCard
            title="Blocks Found"
            value={stats?.blocks_found.toString() || '0'}
            icon="📦"
            color="yellow"
          />
        </div>

        {/* Merged Mining Status */}
        {stats?.merged_mining && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-8">
            <h2 className="text-xl font-bold mb-4">Merged Mining Status</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {stats.chains.map((chain) => (
                <div key={chain} className="bg-gray-50 p-4 rounded-lg">
                  <h3 className="font-semibold text-lg capitalize">{chain}</h3>
                  <p className="text-sm text-gray-600">
                    Network Difficulty: {formatNumber(stats.network_difficulty[chain])}
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Sync Status */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          <h2 className="text-xl font-bold mb-4">Blockchain Sync Status</h2>
          <div className="space-y-4">
            {stats?.sync_status.map((sync) => (
              <div key={sync.chain} className="border-b pb-4 last:border-b-0">
                <div className="flex justify-between items-center mb-2">
                  <h3 className="font-semibold capitalize">{sync.chain}</h3>
                  <span className={`px-3 py-1 rounded-full text-sm ${
                    sync.is_syncing ? 'bg-yellow-200 text-yellow-800' : 'bg-green-200 text-green-800'
                  }`}>
                    {sync.is_syncing ? 'Syncing...' : 'Synced'}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full"
                    style={{ width: `${sync.sync_progress}%` }}
                  ></div>
                </div>
                <div className="flex justify-between text-sm text-gray-600 mt-1">
                  <span>Blocks: {sync.blocks.toLocaleString()}</span>
                  <span>{sync.sync_progress.toFixed(2)}%</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Blocks */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold mb-4">Recent Blocks</h2>
          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Height
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Chain
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Hash
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Time
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {blocks.map((block) => (
                  <tr key={`${block.chain}-${block.height}`}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {block.height}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 capitalize">
                      {block.chain}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                      {block.hash.substring(0, 16)}...
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        block.status === 'confirmed'
                          ? 'bg-green-100 text-green-800'
                          : 'bg-yellow-100 text-yellow-800'
                      }`}>
                        {block.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(block.time).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  );
}

function StatCard({ title, value, icon, color }: { title: string; value: string; icon: string; color: string }) {
  const colors = {
    blue: 'bg-blue-500',
    green: 'bg-green-500',
    purple: 'bg-purple-500',
    yellow: 'bg-yellow-500',
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-gray-500 text-sm">{title}</p>
          <p className="text-2xl font-bold mt-1">{value}</p>
        </div>
        <div className={`${colors[color as keyof typeof colors]} w-12 h-12 rounded-full flex items-center justify-center text-2xl`}>
          {icon}
        </div>
      </div>
    </div>
  );
}

function MinerPage() {
  const [address, setAddress] = useState('');
  const [minerStats, setMinerStats] = useState<MinerStats | null>(null);
  const [loading, setLoading] = useState(false);

  const searchMiner = async () => {
    if (!address) return;
    setLoading(true);
    try {
      const res = await axios.get(`${API_URL}/miner/${address}/stats`);
      setMinerStats(res.data);
    } catch (error) {
      console.error('Error fetching miner stats:', error);
    }
    setLoading(false);
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="container mx-auto">
        <h1 className="text-3xl font-bold mb-6">Miner Statistics</h1>

        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex gap-4">
            <input
              type="text"
              placeholder="Enter your wallet address..."
              className="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && searchMiner()}
            />
            <button
              onClick={searchMiner}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              disabled={loading}
            >
              {loading ? 'Loading...' : 'Search'}
            </button>
          </div>
        </div>

        {minerStats && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <StatCard
                title="Hashrate"
                value={formatHashrate(minerStats.hashrate)}
                icon="⚡"
                color="blue"
              />
              <StatCard
                title="Valid Shares"
                value={minerStats.shares_valid.toString()}
                icon="✓"
                color="green"
              />
              <StatCard
                title="Invalid Shares"
                value={minerStats.shares_invalid.toString()}
                icon="✗"
                color="purple"
              />
            </div>

            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-xl font-bold mb-4">Balance</h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {Object.entries(minerStats.balance).map(([chain, balance]) => (
                  <div key={chain} className="bg-gray-50 p-4 rounded-lg">
                    <p className="text-sm text-gray-600 capitalize">{chain}</p>
                    <p className="text-2xl font-bold">{balance.toFixed(8)}</p>
                  </div>
                ))}
              </div>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-xl font-bold mb-4">Workers</h2>
              <div className="overflow-x-auto">
                <table className="min-w-full">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Worker
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Hashrate
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Valid Shares
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Last Share
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {minerStats.workers.map((worker) => (
                      <tr key={worker.name}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {worker.name}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {formatHashrate(worker.hashrate)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {worker.shares_valid}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {new Date(worker.last_share).toLocaleString()}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function BlocksPage() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedChain, setSelectedChain] = useState<string>('all');
  const [chains, setChains] = useState<string[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [blocksRes, chainsRes] = await Promise.all([
          axios.get(`${API_URL}/pool/blocks?limit=50`),
          axios.get(`${API_URL}/config/chains`)
        ]);
        setBlocks(blocksRes.data);
        setChains(chainsRes.data.chains || []);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching blocks:', error);
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, []);

  const filteredBlocks = selectedChain === 'all'
    ? blocks
    : blocks.filter(b => b.chain === selectedChain);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-blue-600 text-white shadow-lg">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">Mining Pool - Blocks</h1>
            <div className="flex space-x-6">
              <Link to="/" className="hover:text-blue-200">Dashboard</Link>
              <Link to="/miners" className="hover:text-blue-200">Miners</Link>
              <Link to="/blocks" className="hover:text-blue-200">Blocks</Link>
              <Link to="/payments" className="hover:text-blue-200">Payments</Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="container mx-auto px-6 py-8">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-bold">Block History</h2>
            <select
              className="px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={selectedChain}
              onChange={(e) => setSelectedChain(e.target.value)}
            >
              <option value="all">All Chains</option>
              {chains.map((chain) => (
                <option key={chain} value={chain} className="capitalize">
                  {chain}
                </option>
              ))}
            </select>
          </div>

          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Height
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Chain
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Hash
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Difficulty
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Reward
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Confirmations
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Miner
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Time
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredBlocks.map((block) => (
                  <tr key={`${block.chain}-${block.height}`}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {block.height}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 capitalize">
                      {block.chain}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                      {block.hash.substring(0, 16)}...
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {formatNumber(block.difficulty)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {block.reward.toFixed(2)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {block.confirmations}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        block.status === 'confirmed'
                          ? 'bg-green-100 text-green-800'
                          : block.status === 'pending'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {block.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                      {block.miner ? block.miner.substring(0, 12) + '...' : 'N/A'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(block.time).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  );
}

interface Payment {
  id: number;
  address: string;
  amount: number;
  chain: string;
  txid: string;
  time: string;
  status: string;
}

function PaymentsPage() {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedChain, setSelectedChain] = useState<string>('all');
  const [chains, setChains] = useState<string[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [paymentsRes, chainsRes] = await Promise.all([
          axios.get(`${API_URL}/pool/payments?limit=50`),
          axios.get(`${API_URL}/config/chains`)
        ]);
        setPayments(paymentsRes.data);
        setChains(chainsRes.data.chains || []);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching payments:', error);
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, []);

  const filteredPayments = selectedChain === 'all'
    ? payments
    : payments.filter(p => p.chain === selectedChain);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-blue-600 text-white shadow-lg">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">Mining Pool - Payments</h1>
            <div className="flex space-x-6">
              <Link to="/" className="hover:text-blue-200">Dashboard</Link>
              <Link to="/miners" className="hover:text-blue-200">Miners</Link>
              <Link to="/blocks" className="hover:text-blue-200">Blocks</Link>
              <Link to="/payments" className="hover:text-blue-200">Payments</Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="container mx-auto px-6 py-8">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-bold">Payment History</h2>
            <select
              className="px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={selectedChain}
              onChange={(e) => setSelectedChain(e.target.value)}
            >
              <option value="all">All Chains</option>
              {chains.map((chain) => (
                <option key={chain} value={chain} className="capitalize">
                  {chain}
                </option>
              ))}
            </select>
          </div>

          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Chain
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Address
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Amount
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Transaction
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Time
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredPayments.map((payment) => (
                  <tr key={payment.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {payment.id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 capitalize">
                      {payment.chain}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                      {payment.address.substring(0, 12)}...
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {payment.amount.toFixed(8)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                      {payment.txid ? payment.txid.substring(0, 16) + '...' : 'Pending'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        payment.status === 'completed'
                          ? 'bg-green-100 text-green-800'
                          : payment.status === 'pending'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {payment.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(payment.time).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  );
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/miners" element={<MinerPage />} />
        <Route path="/blocks" element={<BlocksPage />} />
        <Route path="/payments" element={<PaymentsPage />} />
      </Routes>
    </Router>
  );
}

function formatHashrate(hashrate: number): string {
  const units = ['H/s', 'KH/s', 'MH/s', 'GH/s', 'TH/s', 'PH/s'];
  let index = 0;
  let value = hashrate;

  while (value >= 1000 && index < units.length - 1) {
    value /= 1000;
    index++;
  }

  return `${value.toFixed(2)} ${units[index]}`;
}

function formatNumber(num: number): string {
  return num.toLocaleString();
}

export default App;
