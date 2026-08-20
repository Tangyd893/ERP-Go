import type { ReactNode } from "react";
import {
  useForm,
  FormProvider,
  type DefaultValues,
  type FieldValues,
  type UseFormReturn,
} from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import type { ZodType } from "zod";
import { Button } from "./ui/button";

export function ProForm<T extends FieldValues>({
  schema,
  defaultValues,
  onSubmit,
  children,
  submitLabel = "提交",
  submitting,
}: {
  schema: ZodType<T>;
  defaultValues?: DefaultValues<T>;
  onSubmit: (values: T) => void | Promise<void>;
  children: (form: UseFormReturn<T>) => ReactNode;
  submitLabel?: string;
  submitting?: boolean;
}) {
  const form = useForm<T>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(schema as any),
    defaultValues,
  });

  return (
    <FormProvider {...form}>
      <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)} noValidate>
        {children(form)}
        <div className="flex justify-end gap-2">
          <Button type="submit" disabled={submitting || form.formState.isSubmitting}>
            {submitting || form.formState.isSubmitting ? "提交中…" : submitLabel}
          </Button>
        </div>
      </form>
    </FormProvider>
  );
}
