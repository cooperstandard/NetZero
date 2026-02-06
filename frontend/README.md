# NetZero Frontend

A modern React TypeScript frontend for the NetZero expense splitting application.

## Features

- User authentication (login/register) with JWT
- Create and join groups
- View group members
- Create transactions with multiple debtors
- View transaction history
- Track balances between members
- Responsive design with Tailwind CSS

## Tech Stack

- React 18
- TypeScript
- Vite (build tool)
- React Router (routing)
- Axios (HTTP client)
- Tailwind CSS (styling)

## Prerequisites

- Node.js 18+ and npm
- NetZero backend API running on port 8080

## Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:3000`

## Development

The Vite dev server is configured to proxy API requests from `/api` to `http://localhost:8080`. This means:
- Frontend runs on `http://localhost:3000`
- API calls to `/api/v1/*` are proxied to `http://localhost:8080/api/v1/*`

## Building for Production

```bash
npm run build
```

The production build will be in the `dist/` directory.

To preview the production build:
```bash
npm run preview
```

## Project Structure

```
src/
├── components/       # Reusable UI components
│   └── Header.tsx
├── context/         # React context providers
│   └── AuthContext.tsx
├── pages/           # Page components
│   ├── Dashboard.tsx
│   ├── GroupDetail.tsx
│   ├── CreateTransaction.tsx
│   ├── Login.tsx
│   └── Register.tsx
├── services/        # API service layer
│   └── api.ts
├── types/           # TypeScript type definitions
│   └── index.ts
├── App.tsx          # Main app component with routing
├── main.tsx         # Application entry point
└── index.css        # Global styles with Tailwind
```

## Key Features Explained

### Authentication
- JWT-based authentication with automatic token refresh
- Tokens stored in localStorage
- Protected routes redirect to login if not authenticated

### Groups
- Create new groups or join existing ones by name
- View all groups you're a member of
- See group members and their details

### Transactions
- Create transactions with a title, description, creditor, and multiple debtors
- Each debt can have a different amount (dollars and cents)
- View all transactions in a group with detailed debt information
- See who owes whom and the amount

### API Integration
- Centralized API service layer using Axios
- Automatic JWT token injection in requests
- Token refresh on 401 errors
- Type-safe API calls with TypeScript

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## Configuration

The API proxy is configured in `vite.config.ts`. If your backend runs on a different port, update the proxy target:

```typescript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:YOUR_PORT',
        changeOrigin: true,
      },
    },
  },
})
```

## License

This project is part of the NetZero expense splitting application.
