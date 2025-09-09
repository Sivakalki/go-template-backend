import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";
import { useSignup } from "../../hooks/useAuth";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { FieldInfo } from "@/components/FieldInfo";

export const Route = createFileRoute("/auth/signup")({
    component: Signup,
});

function Signup() {
    const signup = useSignup();
    const navigate = useNavigate();

    const form = useForm({
        defaultValues: { email: "", password: "", username: "", confirm_password: "" },
        onSubmit: async ({ value }) => {
            await signup.mutateAsync(value);
        },
    });

    return (
        <Card className="max-w-md mx-auto mt-10 shadow-lg flex items-center justify-center border-2 ">
            <CardHeader>
                <CardTitle className="text-2xl">Signup</CardTitle>
            </CardHeader>
            <CardContent>
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        form.handleSubmit().then(() => {
                            if (!signup.isError) {
                                form.reset();
                                navigate({ to: "/auth/login" });
                            }
                        });
                    }}
                    className="space-y-4"
                >

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
                    <div className="grid gap-2">
                        <Label htmlFor="username">Username</Label>
                        <form.Field
                            name="username"
                            validators={{
                                onChange: ({ value }) =>
                                    /\s/.test(value)
                                        ? "Username should not contain spaces"
                                        : value.length < 3
                                            ? "Username must be at least 3 characters"
                                            : undefined,
                            }}
                            children={(field) => (
                                <>

                                    <Input
                                        id="username"
                                        type="text"
                                        value={field.state.value}
                                        onBlur={field.handleBlur}
                                        placeholder="Enter your username"
                                        onChange={(e) => field.handleChange(e.target.value)}
                                    />
                                    <FieldInfo field={field} />
                                </>
                            )}
                        />

                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="password">Password</Label>
                        <form.Field
                            name="password"
                            validators={{
                                onChange: ({ value }) =>
                                    value.length < 8
                                        ? "Password must be at least 8 characters"
                                        : !/[A-Z]/.test(value)
                                            ? "Password must contain at least one uppercase letter"
                                            : !/[a-z]/.test(value)
                                                ? "Password must contain at least one lowercase letter"
                                                : !/[0-9]/.test(value)
                                                    ? "Password must contain at least one number"
                                                    : !/[^A-Za-z0-9]/.test(value)
                                                        ? "Password must contain at least one special character"
                                                        : undefined,

                            }}
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
                    <div className="grid gap-2">
                        <Label htmlFor="confirm_password">Confirm Password</Label>
                        <form.Field
                            name="confirm_password"
                            validators={{
                                onChange: ({ value, fieldApi }) => {
                                    value != fieldApi.form.state.values.password ? "Password do not match" : undefined
                                }
                            }}
                            children={(field) => (
                                <>

                                    <Input
                                        id="confirm_password"
                                        type="confirm_password"
                                        value={field.state.value}
                                        onBlur={field.handleBlur}
                                        placeholder="confirm your password"
                                        onChange={(e) => field.handleChange(e.target.value)}
                                    />
                                    <FieldInfo field={field} />
                                </>
                            )}
                        />

                    </div>

                    <Button type="submit" className="w-full" disabled={signup.isPending}>
                        {signup.isPending ? "Signing up..." : "Signup"}
                    </Button>

                    {signup.isError && (
                        <p className="text-red-600 text-sm mt-2">
                            {(signup.error as Error).message}
                        </p>
                    )}
                </form>
            </CardContent>
        </Card>
    )
}
