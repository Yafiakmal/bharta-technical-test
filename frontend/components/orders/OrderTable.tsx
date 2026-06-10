"use client";

import { Order, Product } from "@/types";
import { formatRupiah } from "@/lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

interface OrderTableProps {
  orders: Order[];
  products: Product[]; // Mungkin masih diperlukan untuk validasi atau keperluan lain
}

export function OrderTable({ orders, products }: OrderTableProps) {
  if (orders.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed py-16 text-center">
        <p className="text-sm font-medium text-muted-foreground">
          No orders yet
        </p>
        <p className="mt-1 text-xs text-muted-foreground">
          Place your first order using the button above.
        </p>
      </div>
    );
  }

  return (
    <div className="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Product</TableHead>
            <TableHead>Price at Order</TableHead>
            <TableHead>Quantity</TableHead>
            <TableHead>Total</TableHead>
            <TableHead>Date</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {orders.map((order) => (
            <TableRow key={order.id}>
              <TableCell className="font-medium">
                {order.product_name}
              </TableCell>
              <TableCell>{formatRupiah(order.product_price)}</TableCell>
              <TableCell>
                <Badge variant="secondary">{order.quantity} unit</Badge>
              </TableCell>
              <TableCell className="font-semibold">
                {formatRupiah(order.product_price * order.quantity)}
              </TableCell>
              <TableCell className="text-muted-foreground">
                {new Date(order.created_at).toLocaleDateString("id-ID", {
                  year: "numeric",
                  month: "long",
                  day: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
