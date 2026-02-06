import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { groupAPI, transactionAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';
import type { GroupMember } from '../types';

interface DebtItem {
  debtor: string;
  dollars: string;
  cents: string;
}

const CreateTransaction = () => {
  const { groupId } = useParams<{ groupId: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [creditor, setCreditor] = useState('');
  const [debts, setDebts] = useState<DebtItem[]>([
    { debtor: '', dollars: '', cents: '0' },
  ]);

  useEffect(() => {
    if (groupId) {
      fetchMembers();
    }
  }, [groupId]);

  useEffect(() => {
    if (user) {
      setCreditor(user.id);
    }
  }, [user]);

  const fetchMembers = async () => {
    try {
      setLoading(true);
      const data = await groupAPI.getMembers(groupId!);
      setMembers(data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load members');
    } finally {
      setLoading(false);
    }
  };

  const addDebtRow = () => {
    setDebts([...debts, { debtor: '', dollars: '', cents: '0' }]);
  };

  const removeDebtRow = (index: number) => {
    if (debts.length > 1) {
      setDebts(debts.filter((_, i) => i !== index));
    }
  };

  const updateDebt = (
    index: number,
    field: keyof DebtItem,
    value: string
  ) => {
    const newDebts = [...debts];
    newDebts[index][field] = value;
    setDebts(newDebts);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!creditor) {
      setError('Please select a creditor');
      return;
    }

    const validDebts = debts.filter(
      (debt) => debt.debtor && (debt.dollars || debt.cents !== '0')
    );

    if (validDebts.length === 0) {
      setError('Please add at least one debt');
      return;
    }

    setSubmitting(true);

    try {
      const transactionData = {
        title,
        description: description || undefined,
        creditor,
        group_id: groupId!,
        transactions: validDebts.map((debt) => ({
          debtor: debt.debtor,
          amount: {
            dollars: parseInt(debt.dollars) || 0,
            cents: parseInt(debt.cents) || 0,
          },
        })),
      };

      await transactionAPI.create(transactionData);
      navigate(`/groups/${groupId}`);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create transaction');
    } finally {
      setSubmitting(false);
    }
  };

  const getTotalAmount = () => {
    return debts.reduce((sum, debt) => {
      const dollars = parseInt(debt.dollars) || 0;
      const cents = parseInt(debt.cents) || 0;
      return sum + dollars + cents / 100;
    }, 0);
  };

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">Loading...</div>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <Link
          to={`/groups/${groupId}`}
          className="text-blue-600 hover:text-blue-800 mb-2 inline-block"
        >
          ← Back to Group
        </Link>
        <h2 className="text-3xl font-bold text-gray-900">New Transaction</h2>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4 mb-6">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6">
        <div className="space-y-6">
          <div>
            <label
              htmlFor="title"
              className="block text-sm font-medium text-gray-700"
            >
              Title *
            </label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="e.g., Dinner at restaurant"
            />
          </div>

          <div>
            <label
              htmlFor="description"
              className="block text-sm font-medium text-gray-700"
            >
              Description
            </label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Optional description"
            />
          </div>

          <div>
            <label
              htmlFor="creditor"
              className="block text-sm font-medium text-gray-700"
            >
              Who paid? *
            </label>
            <select
              id="creditor"
              value={creditor}
              onChange={(e) => setCreditor(e.target.value)}
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select creditor</option>
              {members.map((member) => (
                <option key={member.id} value={member.id}>
                  {member.name || member.email}
                </option>
              ))}
            </select>
          </div>

          <div>
            <div className="flex justify-between items-center mb-3">
              <label className="block text-sm font-medium text-gray-700">
                Who owes money? *
              </label>
              <button
                type="button"
                onClick={addDebtRow}
                className="text-sm text-blue-600 hover:text-blue-800"
              >
                + Add person
              </button>
            </div>

            <div className="space-y-3">
              {debts.map((debt, index) => (
                <div
                  key={index}
                  className="flex items-center space-x-3 p-3 border border-gray-200 rounded-md"
                >
                  <select
                    value={debt.debtor}
                    onChange={(e) =>
                      updateDebt(index, 'debtor', e.target.value)
                    }
                    required
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="">Select person</option>
                    {members.map((member) => (
                      <option key={member.id} value={member.id}>
                        {member.name || member.email}
                      </option>
                    ))}
                  </select>

                  <div className="flex items-center space-x-1">
                    <span className="text-gray-500">$</span>
                    <input
                      type="number"
                      value={debt.dollars}
                      onChange={(e) =>
                        updateDebt(index, 'dollars', e.target.value)
                      }
                      placeholder="0"
                      min="0"
                      className="w-20 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    />
                    <span className="text-gray-500">.</span>
                    <input
                      type="number"
                      value={debt.cents}
                      onChange={(e) =>
                        updateDebt(index, 'cents', e.target.value)
                      }
                      placeholder="00"
                      min="0"
                      max="99"
                      className="w-16 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>

                  {debts.length > 1 && (
                    <button
                      type="button"
                      onClick={() => removeDebtRow(index)}
                      className="text-red-600 hover:text-red-800"
                    >
                      ✕
                    </button>
                  )}
                </div>
              ))}
            </div>

            <div className="mt-4 p-3 bg-gray-50 rounded-md">
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium text-gray-700">
                  Total Amount:
                </span>
                <span className="text-lg font-bold text-gray-900">
                  ${getTotalAmount().toFixed(2)}
                </span>
              </div>
            </div>
          </div>

          <div className="flex justify-end space-x-3 pt-4">
            <Link
              to={`/groups/${groupId}`}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 disabled:bg-green-400"
            >
              {submitting ? 'Creating...' : 'Create Transaction'}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default CreateTransaction;
