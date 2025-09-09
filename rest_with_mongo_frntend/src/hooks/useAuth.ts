import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query"
import { loginRequest, signupRequest, logoutRequest, sessionRequest } from "@/api/auth"


export function useSignup() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: signupRequest,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["session"] });
        }
    })
}


export function useLogin() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: loginRequest,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['session'] })
        }
    })
}

export function useLogout() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: logoutRequest,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['session'] })
        }
    })
}


export function useSession() {
    return useQuery({
        queryKey: ['session'],
        queryFn: sessionRequest,
        retry: false,
    })
}