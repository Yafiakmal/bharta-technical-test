"use client";

import { useCallback, useEffect, useState } from "react";
import { Plus, RefreshCw } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";

import { productApi, orderApi } from "@/lib/api";
import { Product, Order } from "@/types";
import { ProductTable } from "@/components/products/ProductTable";
import {
  ProductForm,
  ProductFormValues,
} from "@/components/products/ProductForm";
import { OrderTable } from "@/components/orders/OrderTable";
import { OrderForm, OrderFormValues } from "@/components/orders/OrderForm";

export default function Home() {
  // Pastikan selalu array, bukan null
  const [products, setProducts] = useState<Product[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [loadingProducts, setLoadingProducts] = useState(true);
  const [loadingOrders, setLoadingOrders] = useState(true);

  const [addProductOpen, setAddProductOpen] = useState(false);
  const [isSubmittingProduct, setIsSubmittingProduct] = useState(false);

  const [addOrderOpen, setAddOrderOpen] = useState(false);
  const [isSubmittingOrder, setIsSubmittingOrder] = useState(false);

  const fetchProducts = useCallback(async () => {
    setLoadingProducts(true);
    try {
      const data = await productApi.getAll();
      setProducts(Array.isArray(data) ? data : []); // Pastikan array
    } catch (err) {
      toast.error((err as Error).message);
      setProducts([]); // Set array kosong jika error
    } finally {
      setLoadingProducts(false);
    }
  }, []);

  const fetchOrders = useCallback(async () => {
    setLoadingOrders(true);
    try {
      const data = await orderApi.getAll();
      setOrders(Array.isArray(data) ? data : []); // Pastikan array
    } catch (err) {
      toast.error((err as Error).message);
      setOrders([]); // Set array kosong jika error
    } finally {
      setLoadingOrders(false);
    }
  }, []);

  useEffect(() => {
    fetchProducts();
    fetchOrders();
  }, [fetchProducts, fetchOrders]);

  async function handleAddProduct(values: ProductFormValues) {
    setIsSubmittingProduct(true);
    try {
      await productApi.create(values);
      toast.success("Product added successfully.");
      setAddProductOpen(false);
      fetchProducts();
    } catch (err) {
      toast.error((err as Error).message);
    } finally {
      setIsSubmittingProduct(false);
    }
  }

  async function handleCreateOrder(values: OrderFormValues) {
    setIsSubmittingOrder(true);
    try {
      await orderApi.create(values);
      toast.success("Order placed successfully.");
      setAddOrderOpen(false);
      fetchOrders();
      fetchProducts();
    } catch (err) {
      toast.error((err as Error).message);
    } finally {
      setIsSubmittingOrder(false);
    }
  }

  return (
    <main className="mx-auto max-w-5xl px-4 py-10">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold tracking-tight">
          Product Management
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage your products and orders in one place.
        </p>
      </div>

      <Tabs defaultValue="products">
        <TabsList className="mb-6">
          <TabsTrigger value="products">Products</TabsTrigger>
          <TabsTrigger value="orders">Orders</TabsTrigger>
        </TabsList>

        {/* ── Products Tab ── */}
        <TabsContent value="products">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-medium">
              {!loadingProducts && products && (
                <span className="text-muted-foreground font-normal text-sm ml-1">
                  {products.length} product{products.length !== 1 ? "s" : ""}
                </span>
              )}
            </h2>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="icon"
                onClick={fetchProducts}
                disabled={loadingProducts}
              >
                <RefreshCw
                  className={`h-4 w-4 ${loadingProducts ? "animate-spin" : ""}`}
                />
                <span className="sr-only">Refresh</span>
              </Button>
              <Dialog open={addProductOpen} onOpenChange={setAddProductOpen}>
                <DialogTrigger>
                  <Button>
                    <Plus className="mr-2 h-4 w-4" />
                    Add Product
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <ProductForm
                    onSubmit={handleAddProduct}
                    isSubmitting={isSubmittingProduct}
                  />
                </DialogContent>
              </Dialog>
            </div>
          </div>

          {loadingProducts ? (
            <TableSkeleton cols={4} />
          ) : (
            <ProductTable products={products || []} onRefresh={fetchProducts} />
          )}
        </TabsContent>

        {/* ── Orders Tab ── */}
        <TabsContent value="orders">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-medium">
              {!loadingOrders && orders && (
                <span className="text-muted-foreground font-normal text-sm ml-1">
                  {orders.length} order{orders.length !== 1 ? "s" : ""}
                </span>
              )}
            </h2>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="icon"
                onClick={fetchOrders}
                disabled={loadingOrders}
              >
                <RefreshCw
                  className={`h-4 w-4 ${loadingOrders ? "animate-spin" : ""}`}
                />
                <span className="sr-only">Refresh</span>
              </Button>
              <Dialog open={addOrderOpen} onOpenChange={setAddOrderOpen}>
                <DialogTrigger>
                  <Button>
                    <Plus className="mr-2 h-4 w-4" />
                    New Order
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <OrderForm
                    products={products || []}
                    onSubmit={handleCreateOrder}
                    isSubmitting={isSubmittingOrder}
                  />
                </DialogContent>
              </Dialog>
            </div>
          </div>

          {loadingOrders ? (
            <TableSkeleton cols={4} />
          ) : (
            <OrderTable orders={orders || []} products={products || []} />
          )}
        </TabsContent>
      </Tabs>
    </main>
  );
}

function TableSkeleton({ cols }: { cols: number }) {
  return (
    <div className="rounded-lg border">
      <div className="p-4 space-y-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="flex gap-4">
            {Array.from({ length: cols }).map((_, j) => (
              <Skeleton key={j} className="h-5 flex-1" />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
}
