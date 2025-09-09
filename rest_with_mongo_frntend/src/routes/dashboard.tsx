import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/dashboard")({
    component: Dashboard,
});

function Dashboard() {
    return (
        <div>
            <h1>Welcome to docs</h1>
        </div>
    )
}