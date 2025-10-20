// web/src/services/api/client.ts
export const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
});

// web/src/store/useAuthStore.ts (Zustand)
export const useAuthStore = create<AuthState>((set, get) => ({
    user: null,
    login: async (credentials) => {
        const user = await authService.login(credentials);
        set({ user });
    }
}));

// web/src/hooks/useWebSocket.ts
export const useWebSocket = (url: string) => {
    // Real-time updates for agent status
};
