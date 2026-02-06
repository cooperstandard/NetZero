import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { groupAPI, transactionAPI } from '../services/api';
import type { GroupMember, Transaction, Debt } from '../types';

const GroupDetail = () => {
  const { groupId } = useParams<{ groupId: string }>();
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [transactionDetails, setTransactionDetails] = useState<{
    [key: string]: Debt[];
  }>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (groupId) {
      fetchGroupData();
    }
  }, [groupId]);

  const fetchGroupData = async () => {
    try {
      setLoading(true);
      const [membersData, transactionsData] = await Promise.all([
        groupAPI.getMembers(groupId!),
        transactionAPI.getTransactions({ group_id: groupId }),
      ]);

      setMembers(membersData);
      setTransactions(transactionsData);

      if (transactionsData.length > 0) {
        const transactionIds = transactionsData.map((t) => t.id);
        const details = await transactionAPI.getTransactionDetails(
          transactionIds
        );
        setTransactionDetails(details);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load group data');
    } finally {
      setLoading(false);
    }
  };

  const getMemberName = (userId: string) => {
    const member = members.find((m) => m.id === userId);
    return member?.name || member?.email || 'Unknown';
  };

  const formatAmount = (amount: string) => {
    return `$${parseFloat(amount).toFixed(2)}`;
  };

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">Loading...</div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex justify-between items-center mb-8">
        <div>
          <Link
            to="/dashboard"
            className="text-blue-600 hover:text-blue-800 mb-2 inline-block"
          >
            ← Back to Dashboard
          </Link>
          <h2 className="text-3xl font-bold text-gray-900">Group Details</h2>
        </div>
        <Link
          to={`/groups/${groupId}/transaction/new`}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
        >
          New Transaction
        </Link>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4 mb-6">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Members ({members.length})
            </h3>
            <ul className="space-y-3">
              {members.map((member) => (
                <li key={member.id} className="flex items-center">
                  <div className="flex-shrink-0 h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                    {member.name?.[0]?.toUpperCase() ||
                      member.email[0].toUpperCase()}
                  </div>
                  <div className="ml-3">
                    <p className="text-sm font-medium text-gray-900">
                      {member.name || 'No name'}
                    </p>
                    <p className="text-xs text-gray-500">{member.email}</p>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="lg:col-span-2">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Transactions
            </h3>
            {transactions.length === 0 ? (
              <p className="text-gray-500 text-center py-8">
                No transactions yet. Create your first transaction to get
                started.
              </p>
            ) : (
              <div className="space-y-4">
                {transactions.map((transaction) => {
                  const debts = transactionDetails[transaction.id] || [];
                  const totalAmount = debts.reduce(
                    (sum, debt) => sum + parseFloat(debt.amount),
                    0
                  );

                  return (
                    <div
                      key={transaction.id}
                      className="border border-gray-200 rounded-lg p-4"
                    >
                      <div className="flex justify-between items-start mb-2">
                        <div>
                          <h4 className="text-md font-semibold text-gray-900">
                            {transaction.title}
                          </h4>
                          {transaction.description && (
                            <p className="text-sm text-gray-600 mt-1">
                              {transaction.description}
                            </p>
                          )}
                        </div>
                        <span className="text-lg font-bold text-green-600">
                          {formatAmount(totalAmount.toString())}
                        </span>
                      </div>
                      <div className="text-xs text-gray-500 mb-3">
                        By {getMemberName(transaction.author_id)} on{' '}
                        {new Date(transaction.created_at).toLocaleString()}
                      </div>

                      {debts.length > 0 && (
                        <div className="mt-3 border-t border-gray-100 pt-3">
                          <p className="text-xs font-semibold text-gray-700 mb-2">
                            Debts:
                          </p>
                          <div className="space-y-1">
                            {debts.map((debt) => (
                              <div
                                key={debt.id}
                                className="flex justify-between items-center text-sm"
                              >
                                <span className="text-gray-700">
                                  {getMemberName(debt.debtor)} owes{' '}
                                  {getMemberName(debt.creditor)}
                                </span>
                                <span
                                  className={`font-semibold ${
                                    debt.paid
                                      ? 'text-green-600'
                                      : 'text-orange-600'
                                  }`}
                                >
                                  {formatAmount(debt.amount)}
                                  {debt.paid && ' ✓'}
                                </span>
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default GroupDetail;
