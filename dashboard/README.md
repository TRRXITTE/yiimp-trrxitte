# Mining Pool Dashboard

Modern web dashboard for the multicoin merged mining pool, built with React, TypeScript, and Tailwind CSS.

## Features

- Real-time pool statistics monitoring
- Multi-chain support with merged mining status display
- Blockchain sync status with progress indicators
- Block history with chain filtering
- Payment history tracking
- Miner statistics lookup
- Worker performance monitoring
- Responsive design with mobile support

## Prerequisites

- Node.js 18+ and npm
- Mining pool API server running on port 8001

## Installation

```bash
cd dashboard
npm install
```

## Development

Start the development server:

```bash
npm run dev
```

The dashboard will be available at `http://localhost:3000`

## Building for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

## Configuration

The dashboard connects to the pool API at `http://localhost:8001/api` by default. To change this, edit the `API_URL` constant in `src/App.tsx`.

## Pages

- **Dashboard** (`/`) - Pool overview with stats, sync status, and recent blocks
- **Miners** (`/miners`) - Search and view miner statistics
- **Blocks** (`/blocks`) - Block history with filtering by chain
- **Payments** (`/payments`) - Payment history with filtering by chain

## Technology Stack

- React 18
- TypeScript 5
- Vite 5 (build tool)
- Tailwind CSS 3 (styling)
- React Router 6 (navigation)
- Recharts (charts and visualizations)
- Axios (HTTP client)

## Auto-refresh

The dashboard automatically refreshes data:
- Pool stats: every 10 seconds
- Blocks: every 30 seconds
- Payments: every 30 seconds
