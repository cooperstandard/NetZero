export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface Group {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: string;
  title: string;
  description: string;
  author_id: string;
  group_id: string;
  created_at: string;
  updated_at: string;
}

export interface Debt {
  id: string;
  amount: string;
  transaction_id: string;
  debtor: string;
  creditor: string;
  created_at: string;
  updated_at: string;
  paid: boolean;
}

export interface Balance {
  user_id: string;
  group_id: string;
  creditor_id: string;
  balance: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
  expiresInSeconds?: number;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface LoginResponse {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
  token: string;
  refresh_token: string;
}

export interface RegisterResponse {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTransactionRequest {
  title: string;
  description?: string;
  creditor: string;
  group_id: string;
  transaction_id?: string;
  transactions: {
    debtor: string;
    amount: {
      dollars: number;
      cents: number;
    };
  }[];
}

export interface CreateTransactionResponse {
  transaction_id: string;
  transactions: Debt[];
}

export interface GroupMember {
  id: string;
  name: string;
  email: string;
}
