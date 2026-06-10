"use client";

import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Product } from "@/types";

// Perbaikan schema untuk quantity
const formSchema = z.object({
  productId: z.string().min(1, "Please select a product"),
  quantity: z.number().int().positive("Quantity must be at least 1"),
});

export type OrderFormValues = z.infer<typeof formSchema>;

interface OrderFormProps {
  products: Product[];
  onSubmit: (values: OrderFormValues) => Promise<void>;
  isSubmitting: boolean;
}

export function OrderForm({
  products,
  onSubmit,
  isSubmitting,
}: OrderFormProps) {
  const form = useForm<OrderFormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      productId: "",
      quantity: 1,
    },
  });

  const selectedProductId = form.watch("productId");
  const selectedProduct = products.find((p) => p.id === selectedProductId);

  const handleSubmit = (values: OrderFormValues) => {
    return onSubmit(values);
  };

  return (
    <>
      <DialogHeader>
        <DialogTitle>Create Order</DialogTitle>
        <DialogDescription>
          Select a product and enter the quantity to place an order.
        </DialogDescription>
      </DialogHeader>

      <form
        id="order-form"
        onSubmit={form.handleSubmit(handleSubmit)}
        className="pt-2"
      >
        <FieldGroup>
          <Controller
            name="productId"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="order-product">Product</FieldLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <SelectTrigger
                    id="order-product"
                    aria-invalid={fieldState.invalid}
                  >
                    <SelectValue placeholder="Select a product" />
                  </SelectTrigger>
                  <SelectContent>
                    {products
                      .filter((p) => p.stock > 0)
                      .map((product) => (
                        <SelectItem key={product.id} value={product.id}>
                          {product.name}
                          <span className="ml-2 text-muted-foreground text-xs">
                            ({product.stock} available)
                          </span>
                        </SelectItem>
                      ))}
                    {products.filter((p) => p.stock > 0).length === 0 && (
                      <div className="px-2 py-4 text-center text-sm text-muted-foreground">
                        No products in stock
                      </div>
                    )}
                  </SelectContent>
                </Select>
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />

          <Controller
            name="quantity"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="order-quantity">Quantity</FieldLabel>
                <Input
                  id="order-quantity"
                  type="number"
                  min={1}
                  max={selectedProduct?.stock ?? undefined}
                  step={1}
                  aria-invalid={fieldState.invalid}
                  onChange={(e) => field.onChange(parseInt(e.target.value))}
                  value={field.value || ""}
                />
                {selectedProduct && (
                  <p className="text-xs text-muted-foreground mt-1">
                    Max: {selectedProduct.stock} unit
                  </p>
                )}
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
        </FieldGroup>

        <div className="flex justify-end pt-4">
          <Button type="submit" disabled={isSubmitting} form="order-form">
            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Place Order
          </Button>
        </div>
      </form>
    </>
  );
}
