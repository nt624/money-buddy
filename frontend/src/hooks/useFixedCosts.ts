import { useEffect, useState } from "react";
import {
  createFixedCost,
  getFixedCosts,
  updateFixedCost,
  deleteFixedCost,
} from "@/lib/api/fixed-costs";
import { FixedCost, CreateFixedCostInput, UpdateFixedCostInput } from "@/lib/types/fixed-cost";

export function useFixedCosts() {
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
        const data = await getFixedCosts();
        setFixedCosts(data.fixed_costs);
      } catch (err) {
        setError(err instanceof Error ? err.message : "unknown error");
      } finally {
        setIsLoading(false);
      }
    };

    fetchFixedCosts();
  }, []);

  // Refetch
  const refetchFixedCosts = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await getFixedCosts();
      setFixedCosts(data.fixed_costs);
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown error");
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
      const newFixedCost = await createFixedCost(input);
      setFixedCosts((prevFixedCosts) => [...prevFixedCosts, newFixedCost]);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown error");
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
      const updatedFixedCost = await updateFixedCost(id, input);
      setFixedCosts((prevFixedCosts) =>
        prevFixedCosts.map((fc) => (fc.id === id ? updatedFixedCost : fc))
      );
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown error");
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
      await deleteFixedCost(id);
      setFixedCosts((prevFixedCosts) =>
        prevFixedCosts.filter((fc) => fc.id !== id)
      );
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown error");
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
