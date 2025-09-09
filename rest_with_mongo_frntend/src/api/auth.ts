import api from "./client";

export type AuthResponse = {
    user: { id: string; email: string };
    message?: string;
}


// signup
export const signupRequest = async (data: { email: string; username: string; password: string; confirm_password: string }): Promise<AuthResponse> => {
    const res = await api.post("/auth/register", data, { withCredentials: true });
    return res.data
}

// login
export const loginRequest = async (data: { email: string, password: string }): Promise<AuthResponse> => {
    const res = await api.post("/auth/login", data, { withCredentials: true });
    return res.data
}


//logout
export const logoutRequest = async (): Promise<{ message: string }> => {
    const res = await api.post("/auth/logout", {}, { withCredentials: true });
    return res.data;
}

export const sessionRequest = async (): Promise<AuthResponse> => {
    const res = await api.get("/auth/session", { withCredentials: true });
    return res.data;
};