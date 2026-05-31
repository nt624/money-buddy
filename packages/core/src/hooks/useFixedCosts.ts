import { useEffect, useState } from "react";
import { useDataSource } from "../data";
import { FixedCost, CreateFixedCostInput, UpdateFixedCostInput } from "../types/fixed-cost";

export function useFixedCosts() {
  const ds = useDataSource();
  const [fixedCosts, setFixedCosts] = useState<FixedCost[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // GET
  useEffect(() => {
    const fetchFixedCosts = async () => {
      setIsLoading(true);
      setError(null);

      try {
        const data = await ds.fixedCosts.list();
        setFixedCosts(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "エラーが発生しました");
      } finally {
        setIsLoading(false);
      }
    };

    fetchFixedCosts();
  }, [ds]);

  // Refetch
  const refetchFixedCosts = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await ds.fixedCosts.list();
      setFixedCosts(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  };

  // POST
  const handleCreateFixedCost = async (
    input: CreateFixedCostInput
  ): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);

    try {
      const newFixedCost = await ds.fixedCosts.create(input);
      setFixedCosts((prevFixedCosts) => [...prevFixedCosts, newFixedCost]);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  // PUT
  const handleUpdateFixedCost = async (
    id: number,
    input: UpdateFixedCostInput
  ): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);

    try {
      const updatedFixedCost = await ds.fixedCosts.update(id, input);
      setFixedCosts((prevFixedCosts) =>
        prevFixedCosts.map((fc) => (fc.id === id ? updatedFixedCost : fc))
      );
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  // DELETE
  const handleDeleteFixedCost = async (id: number): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);

    try {
      await ds.fixedCosts.remove(id);
      setFixedCosts((prevFixedCosts) =>
        prevFixedCosts.filter((fc) => fc.id !== id)
      );
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    fixedCosts,
    isLoading,
    isSubmitting,
    error,
    createFixedCost: handleCreateFixedCost,
    updateFixedCost: handleUpdateFixedCost,
    deleteFixedCost: handleDeleteFixedCost,
    refetch: refetchFixedCosts,
  };
}
