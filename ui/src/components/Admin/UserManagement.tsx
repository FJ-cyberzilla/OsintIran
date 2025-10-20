// ui/src/components/Admin/UserManagement.tsx
import React, { useState, useEffect } from 'react';
import { useAdminAPI } from '../../../hooks/useAdminAPI';

export const UserManagement: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'users' | 'tenants' | 'api-keys'>('users');
  const [users, setUsers] = useState<User[]>([]);
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const { getUsers, getTenants, getAPIKeys, createUser, createTenant, createAPIKey } = useAdminAPI();

  useEffect(() => {
    loadData();
  }, [activeTab]);

  const loadData = async () => {
    setIsLoading(true);
    try {
      switch (activeTab) {
        case 'users':
          const usersData = await getUsers();
          setUsers(usersData);
          break;
        case 'tenants':
          const tenantsData = await getTenants();
          setTenants(tenantsData);
          break;
        case 'api-keys':
          const apiKeysData = await getAPIKeys();
          setApiKeys(apiKeysData);
          break;
      }
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="user-management">
      <div className="management-header">
        <h1>👥 User & Tenant Management</h1>
        <p>Manage users, tenants, and API keys across the platform</p>
      </div>

      {/* Tab Navigation */}
      <div className="management-tabs">
        <button 
          className={activeTab === 'users' ? 'active' : ''}
          onClick={() => setActiveTab('users')}
        >
          👤 Users
        </button>
        <button 
          className={activeTab === 'tenants' ? 'active' : ''}
          onClick={() => setActiveTab('tenants')}
        >
          🏢 Tenants
        </button>
        <button 
          className={activeTab === 'api-keys' ? 'active' : ''}
          onClick={() => setActiveTab('api-keys')}
        >
          🔑 API Keys
        </button>
      </div>

      {/* Users Tab */}
      {activeTab === 'users' && (
        <UsersTab 
          users={users} 
          isLoading={isLoading}
          onCreateUser={createUser}
          onRefresh={loadData}
        />
      )}

      {/* Tenants Tab */}
      {activeTab === 'tenants' && (
        <TenantsTab 
          tenants={tenants}
          isLoading={isLoading}
          onCreateTenant={createTenant}
          onRefresh={loadData}
        />
      )}

      {/* API Keys Tab */}
      {activeTab === 'api-keys' && (
        <APIKeysTab
          apiKeys={apiKeys}
          isLoading={isLoading}
          onCreateAPIKey={createAPIKey}
          onRefresh={loadData}
        />
      )}
    </div>
  );
};

const UsersTab: React.FC<{
  users: User[];
  isLoading: boolean;
  onCreateUser: (user: CreateUserRequest) => Promise<void>;
  onRefresh: () => void;
}> = ({ users, isLoading, onCreateUser, onRefresh }) => {
  const [showCreateForm, setShowCreateForm] = useState(false);

  return (
    <div className="users-tab">
      <div className="tab-header">
        <h3>User Management</h3>
        <button 
          className="btn-primary"
          onClick={() => setShowCreateForm(true)}
        >
          + Add User
        </button>
      </div>

      {showCreateForm && (
        <CreateUserForm
          onSubmit={onCreateUser}
          onCancel={() => setShowCreateForm(false)}
          onSuccess={() => {
            setShowCreateForm(false);
            onRefresh();
          }}
        />
      )}

      {isLoading ? (
        <div className="loading">Loading users...</div>
      ) : (
        <div className="users-table">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Tenant</th>
                <th>Status</th>
                <th>Last Login</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map(user => (
                <UserRow key={user.id} user={user} />
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

const TenantsTab: React.FC<{
  tenants: Tenant[];
  isLoading: boolean;
  onCreateTenant: (tenant: CreateTenantRequest) => Promise<void>;
  onRefresh: () => void;
}> = ({ tenants, isLoading, onCreateTenant, onRefresh }) => {
  const [showCreateForm, setShowCreateForm] = useState(false);

  return (
    <div className="tenants-tab">
      <div className="tab-header">
        <h3>Tenant Management</h3>
        <button 
          className="btn-primary"
          onClick={() => setShowCreateForm(true)}
        >
          + Create Tenant
        </button>
      </div>

      {showCreateForm && (
        <CreateTenantForm
          onSubmit={onCreateTenant}
          onCancel={() => setShowCreateForm(false)}
          onSuccess={() => {
            setShowCreateForm(false);
            onRefresh();
          }}
        />
      )}

      {isLoading ? (
        <div className="loading">Loading tenants...</div>
      ) : (
        <div className="tenants-grid">
          {tenants.map(tenant => (
            <TenantCard key={tenant.id} tenant={tenant} />
          ))}
        </div>
      )}
    </div>
  );
};

const TenantCard: React.FC<{ tenant: Tenant }> = ({ tenant }) => (
  <div className={`tenant-card plan-${tenant.planType}`}>
    <div className="card-header">
      <h4>{tenant.name}</h4>
      <span className={`status-badge status-${tenant.status}`}>
        {tenant.status}
      </span>
    </div>
    
    <div className="card-body">
      <div className="tenant-info">
        <div className="info-item">
          <span className="label">Plan:</span>
          <span className="value">{tenant.planType}</span>
        </div>
        <div className="info-item">
          <span className="label">Slug:</span>
          <span className="value">{tenant.slug}</span>
        </div>
        <div className="info-item">
          <span className="label">Users:</span>
          <span className="value">{tenant.userCount}/{tenant.maxUsers}</span>
        </div>
        <div className="info-item">
          <span className="label">Requests:</span>
          <span className="value">
            {tenant.usage?.requestsCount.toLocaleString()}/{tenant.maxRequestsPerMonth.toLocaleString()}
          </span>
        </div>
    </div>

    <div className="card-actions">
      <button className="btn-sm">Edit</button>
      <button className="btn-sm">Usage</button>
      <button className="btn-sm">Billing</button>
    </div>
  </div>
);
