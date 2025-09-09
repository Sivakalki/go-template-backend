import { FieldApi } from "@tanstack/react-form";

interface FieldInfoProps {
    field: FieldApi<any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any>;
}

export function FieldInfo({ field }: FieldInfoProps) {
    return (
        <div className="text-sm text-red-500">
            {field.state.meta.errors?.length > 0 && field.state.meta.errors.join(", ")}
        </div>
    );
}
