import { defineStore } from "pinia";

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
}

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => ({
    user: null,
    token: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === "admin",
    currentUser: (state) => state.user,
  },

  actions: {
    // =========================
   // =========================
    // LOGIN (Fixed Execution Flow)
    // =========================
    async login(email: string, password: string) {
      try {
        const response: any = await $fetch("http://localhost:8080/v1/graphql", {
          method: "POST",
          body: {
            query: `
              mutation LoginUser($input: LoginUserInput!) {
                loginUser(input: $input) {
                  token
                  user_id
                  message
                }
              }
            `,
            variables: {
              input: { email, password },
            },
          },
        });

        console.log("Login response:", response);

        if (response?.errors) {
          return {
            success: false,
            error: response.errors[0]?.message || "Login failed",
          };
        }

        const data = response?.data?.loginUser;

        if (!data?.token || !data?.user_id) {
          return {
            success: false,
            error: "Invalid credentials or token missing",
          };
        }

        // Fetch user details using the newly received token explicitly in headers
        const userResponse: any = await $fetch(
          "http://localhost:8080/v1/graphql",
          {
            method: "POST",
            body: {
              query: `
                query GetUser($id: uuid!) {
                  users_by_pk(id: $id) {
                    id
                    email
                    name
                    role
                  }
                }
              `,
              variables: { id: data.user_id },
            },
            headers: {
              Authorization: `Bearer ${data.token}`,
            },
          },
        );

        if (userResponse?.errors) {
          return {
            success: false,
            error: userResponse.errors[0]?.message || "Failed to retrieve profile data",
          };
        }

        const user = userResponse?.data?.users_by_pk;

        if (!user) {
          return { success: false, error: "User profile record not found" };
        }

        // Commit both token and profile to localStorage and application state
        this.setAuth(data.token, user);

        return { success: true, user };
      } catch (error: any) {
        console.error("Login catch block triggered error:", error);
        return {
          success: false,
          error: error?.data?.errors?.[0]?.message || error?.message || "Login failed",
        };
      }
    },

    // =========================
    // REGISTER
    // =========================
    async register(email: string, password: string, name: string) {
      try {
        const response: any = await $fetch("http://localhost:8080/v1/graphql", {
          method: "POST",
          body: {
            query: `
              mutation RegisterUser($input: RegisterUserInput!) {
                registerUser(input: $input) {
                  token
                  user_id
                  message
                }
              }
            `,
            variables: {
              input: { name, email, password },
            },
          },
        });

        console.log("Register response:", response);

        // Check for errors in response
        if (response?.errors) {
          return {
            success: false,
            error: response.errors[0]?.message || "Registration failed",
          };
        }

        const data = response?.data?.registerUser;

        if (!data?.token || !data?.user_id) {
          return {
            success: false,
            error: response?.errors?.[0]?.message || "Registration failed",
          };
        }

        // Fetch user details
        const userResponse: any = await $fetch(
          "http://localhost:8080/v1/graphql",
          {
            method: "POST",
            body: {
              query: `
                query GetUser($id: uuid!) {
                  users_by_pk(id: $id) {
                    id
                    email
                    name
                    role
                  }
                }
              `,
              variables: { id: data.user_id },
            },
            headers: {
              Authorization: `Bearer ${data.token}`,
            },
          },
        );

        if (userResponse?.errors) {
          return {
            success: false,
            error: userResponse.errors[0]?.message || "Failed to fetch user",
          };
        }

        const user = userResponse?.data?.users_by_pk;

        if (!user) {
          return {
            success: false,
            error: "User fetch failed after registration",
          };
        }

        this.setAuth(data.token, user);

        return { success: true, user };
      } catch (error: any) {
        console.error("Registration error:", error);
        return {
          success: false,
          error: error?.message || "Registration failed",
        };
      }
    },

    // =========================
    // AUTH HELPERS
    // =========================
    setAuth(token: string, user: User) {
      this.token = token;
      this.user = user;

      if (process.client) {
        localStorage.setItem("auth_token", token);
        localStorage.setItem("auth_user", JSON.stringify(user));
      }
    },

    logout() {
      this.token = null;
      this.user = null;

      if (process.client) {
        localStorage.removeItem("auth_token");
        localStorage.removeItem("auth_user");
      }

      navigateTo("/login");
    },

    loadAuth() {
      if (!process.client) return;

      const token = localStorage.getItem("auth_token");
      const userStr = localStorage.getItem("auth_user");

      if (token && userStr) {
        try {
          this.token = token;
          this.user = JSON.parse(userStr);
        } catch {
          this.logout();
        }
      }
    },
  },
}); 