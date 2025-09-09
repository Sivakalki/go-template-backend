import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";
import { useLogin, } from "../../hooks/useAuth";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { FieldInfo } from "@/components/FieldInfo";

export const Route = createFileRoute("/auth/login")({
    component: Login,
});

function Login() {
    const login = useLogin();
    const navigate = useNavigate();

    const form = useForm({
        defaultValues: {
            email: "",
            password: "",
        },
        onSubmit: async ({ value }) => {
            try {
                await login.mutateAsync(value);
                navigate({ to: "/dashboard" });
            } catch (error) {
                console.error("Login failed", error);
            }
        },
    });

    return (
        <Card className="max-w-md mx-auto mt-10 shadow-lg">
            <CardHeader>
                <CardTitle className="text-2xl">Login</CardTitle>
            </CardHeader>
            <CardContent>
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        form.handleSubmit().then(() => {
                            if (!login.isError) {
                                form.reset();
                                navigate({ to: "/dashboard" });
                            }
                        });
                    }}
                    className="space-y-4">
                    {/* Email */}
                    <div className="grid gap-2">
                        <Label htmlFor="email">Email</Label>
                        <form.Field
                            name="email"
                            validators={{
                                onChange: ({ value }) =>
                                    !value.includes("@") ? "Invalid email address" : undefined,
                            }}
                            children={(field) => (
                                <>
                                    <Input
                                        id="email"
                                        type="email"
                                        value={field.state.value}
                                        onBlur={field.handleBlur}
                                        placeholder="Enter your email"
                                        onChange={(e) => field.handleChange(e.target.value)}
                                    />
                                    <FieldInfo field={field} />
                                </>
                            )}
                        />
                    </div>

                    {/* Password */}
                    <div className="grid gap-2">
                        <Label htmlFor="password">Password</Label>
                        <form.Field
                            name="password"

                            children={(field) => (
                                <>
                                    <Input
                                        id="password"
                                        type="password"
                                        value={field.state.value}
                                        onBlur={field.handleBlur}
                                        placeholder="Enter your password"
                                        onChange={(e) => field.handleChange(e.target.value)}
                                    />
                                    <FieldInfo field={field} />
                                </>
                            )}
                        />
                    </div>

                    <Button type="submit" className="w-full" disabled={login.isPending}>
                        {login.isPending ? "Logging in..." : "Login"}
                    </Button>

                    {login.isError && (
                        <p className="text-red-600 text-sm mt-2">
                            {(login.error as Error).message}
                        </p>
                    )}
                </form>
            </CardContent>
        </Card>
    );
}
